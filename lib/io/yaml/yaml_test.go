package yaml_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2" // Import the original yaml library to use its Unmarshal/Marshal directly for comparison

	cut "github.com/walterjwhite/go-code/lib/io/yaml" // Code Under Test
)

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Enabled bool   `yaml:"enabled"`
	Items   []struct {
		ID   int    `yaml:"id"`
		Type string `yaml:"type"`
	} `yaml:"items"`
}

func TestRead_Success(t *testing.T) {
	tempDir := t.TempDir()
	yamlContent := `
name: TestApp
version: "1.0.0"
enabled: true
items:
  - id: 1
    type: A
  - id: 2
    type: B
`
	tempFilePath := filepath.Join(tempDir, "test_config.yaml")
	err := os.WriteFile(tempFilePath, []byte(yamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary YAML file")

	var cfg Config
	err = cut.Read(tempFilePath, &cfg)

	require.NoError(t, err, "Read function returned an error")
	assert.Equal(t, "TestApp", cfg.Name)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.True(t, cfg.Enabled)
	require.Len(t, cfg.Items, 2)
	assert.Equal(t, 1, cfg.Items[0].ID)
	assert.Equal(t, "A", cfg.Items[0].Type)
	assert.Equal(t, 2, cfg.Items[1].ID)
	assert.Equal(t, "B", cfg.Items[1].Type)
}

func TestRead_FileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "non_existent.yaml")

	var cfg Config
	err := cut.Read(nonExistentFile, &cfg)

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err), "Error should indicate file not found")
}

func TestRead_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	invalidYamlContent := `
name: TestApp
version: "1.0.0"
enabled: true
items:
  - id: 1
    type: A
  - id: 2
    type: B: Invalid
` // Malformed YAML with an extra colon

	tempFilePath := filepath.Join(tempDir, "invalid_config.yaml")
	err := os.WriteFile(tempFilePath, []byte(invalidYamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary invalid YAML file")

	var cfg Config
	err = cut.Read(tempFilePath, &cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: line 9: mapping values are not allowed in this context")
}

func TestWrite_Success(t *testing.T) {
	cfg := Config{
		Name:    "TestWriteApp",
		Version: "2.0.0",
		Enabled: false,
		Items: []struct {
			ID   int    `yaml:"id"`
			Type string `yaml:"type"`
		}{
			{ID: 3, Type: "C"},
			{ID: 4, Type: "D"},
		},
	}

	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "output_config.yaml")

	err := cut.Write(cfg, outputFilePath)
	require.NoError(t, err, "Write function returned an error")

	writtenBytes, err := os.ReadFile(outputFilePath)
	require.NoError(t, err, "Failed to read written file")

	var readCfg Config
	err = yaml.Unmarshal(writtenBytes, &readCfg) // Use original yaml.Unmarshal for verification
	require.NoError(t, err, "Failed to unmarshal written YAML content for verification")

	assert.Equal(t, cfg.Name, readCfg.Name)
	assert.Equal(t, cfg.Version, readCfg.Version)
	assert.Equal(t, cfg.Enabled, readCfg.Enabled)
	require.Len(t, readCfg.Items, 2)
	assert.Equal(t, cfg.Items[0].ID, readCfg.Items[0].ID)
	assert.Equal(t, cfg.Items[0].Type, readCfg.Items[0].Type)
	assert.Equal(t, cfg.Items[1].ID, readCfg.Items[1].ID)
	assert.Equal(t, cfg.Items[1].Type, readCfg.Items[1].Type)
}

func TestWrite_InvalidInput(t *testing.T) {
	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "invalid_input.yaml")

	unmarshalable := make(chan int)
	err := cut.Write(unmarshalable, outputFilePath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml.Marshal panicked: cannot marshal type: chan int")
}

func TestWrite_UnwritablePath(t *testing.T) {
	tempDir := t.TempDir() // A directory is unwritable as a file

	err := cut.Write(Config{}, tempDir)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is a directory")
}

func TestRead_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	emptyFilePath := filepath.Join(tempDir, "empty.yaml")
	err := os.WriteFile(emptyFilePath, []byte(""), 0644)
	require.NoError(t, err, "Failed to create empty file")

	var cfg Config
	err = cut.Read(emptyFilePath, &cfg)

	require.NoError(t, err)
	assert.Equal(t, Config{}, cfg)
}

func TestRead_IntoNilPointer(t *testing.T) {
	tempDir := t.TempDir()
	yamlContent := `name: TestApp`
	tempFilePath := filepath.Join(tempDir, "nil_ptr.yaml")
	err := os.WriteFile(tempFilePath, []byte(yamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary YAML file")

	var cfg *Config = nil
	err = cut.Read(tempFilePath, cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml.Unmarshal panicked: reflect: reflect.Value.Set using unaddressable value")
}

func TestWrite_EmptyStruct(t *testing.T) {
	emptyCfg := Config{}
	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "empty_struct_output.yaml")

	err := cut.Write(emptyCfg, outputFilePath)
	require.NoError(t, err)

	writtenBytes, err := os.ReadFile(outputFilePath)
	require.NoError(t, err)

	expectedYaml := `name: ""
version: ""
enabled: false
items: []
`
	assert.Equal(t, []byte(expectedYaml), writtenBytes)
}

func TestRead_DifferentStruct(t *testing.T) {
	tempDir := t.TempDir()
	yamlContent := `
name: OriginalName
unused_field: SomeValue
`
	tempFilePath := filepath.Join(tempDir, "diff_struct.yaml")
	err := os.WriteFile(tempFilePath, []byte(yamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary YAML file")

	type SmallConfig struct {
		Name string `yaml:"name"`
	}
	var smallCfg SmallConfig
	err = cut.Read(tempFilePath, &smallCfg)

	require.NoError(t, err)
	assert.Equal(t, "OriginalName", smallCfg.Name)
}

func TestWrite_DifferentStruct(t *testing.T) {
	type CustomConfig struct {
		Title string `yaml:"title"`
		Value int    `yaml:"value"`
	}
	customCfg := CustomConfig{
		Title: "Custom Title",
		Value: 123,
	}

	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "custom_struct_output.yaml")

	err := cut.Write(customCfg, outputFilePath)
	require.NoError(t, err)

	writtenBytes, err := os.ReadFile(outputFilePath)
	require.NoError(t, err)

	var readCustomCfg CustomConfig
	err = yaml.Unmarshal(writtenBytes, &readCustomCfg)
	require.NoError(t, err)

	assert.Equal(t, customCfg.Title, readCustomCfg.Title)
	assert.Equal(t, customCfg.Value, readCustomCfg.Value)

	expectedYaml := "title: Custom Title\nvalue: 123\n"
	assert.Equal(t, expectedYaml, string(writtenBytes))
}

func TestRead_LoggerCoverage(t *testing.T) {

	tempDir := t.TempDir()
	yamlContent := `name: LogTest`
	tempFilePath := filepath.Join(tempDir, "log_read_config.yaml")
	err := os.WriteFile(tempFilePath, []byte(yamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary YAML file")

	var cfg Config
	err = cut.Read(tempFilePath, &cfg)
	require.NoError(t, err)
	assert.Equal(t, "LogTest", cfg.Name)
}

func TestWrite_LoggerCoverage(t *testing.T) {
	cfg := Config{Name: "LogWriteTest"}
	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "log_write_output.yaml")

	err := cut.Write(cfg, outputFilePath)
	require.NoError(t, err)

	_, err = os.Stat(outputFilePath)
	assert.False(t, os.IsNotExist(err), "Output file should exist")
}

type PanicOnUnmarshal struct{}

func (p *PanicOnUnmarshal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	panic("I am a string panic during unmarshal!") // Panic with a string
}

func TestRead_UnmarshalNonErrorPanic(t *testing.T) {
	tempDir := t.TempDir()
	yamlContent := `key: value`
	tempFilePath := filepath.Join(tempDir, "non_error_panic.yaml")
	err := os.WriteFile(tempFilePath, []byte(yamlContent), 0644)
	require.NoError(t, err, "Failed to write temporary YAML file")

	var p PanicOnUnmarshal
	err = cut.Read(tempFilePath, &p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml.Unmarshal panicked: I am a string panic during unmarshal!")
}

type PanicOnMarshal struct{}

func (p PanicOnMarshal) MarshalYAML() (interface{}, error) {
	panic(123) // Panic with an int
}

func TestWrite_MarshalNonErrorPanic(t *testing.T) {
	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "non_error_panic_write.yaml")

	p := PanicOnMarshal{}
	err := cut.Write(p, outputFilePath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml.Marshal panicked: 123")
}
