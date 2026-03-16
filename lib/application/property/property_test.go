package property

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/sflags/gen/gflag"
)

type TestConfig struct {
	Name    string `yaml:"name"`
	Value   int    `yaml:"value"`
	Enabled bool   `yaml:"enabled"`
	Nested  struct {
		Field1 string  `yaml:"field1"`
		Field2 float64 `yaml:"field2"`
	} `yaml:"nested"`
	DeepNested struct {
		Level2 struct {
			DeepField string `yaml:"deepField"`
		} `yaml:"level2"`
	} `yaml:"deepNested"`
	SecretField string `yaml:"SecretField"`
}

func (c *TestConfig) SecretFields() []string {
	return []string{"SecretField", "NonExistentField", "unexportedField", "NonStringField", "EmptySecret", "Nested.Field1", "DeepNested.Level2.DeepField"}
}

type TestConfigWithUnexported struct {
	NonStringField int    `yaml:"NonStringField"`
	EmptySecret    string `yaml:"EmptySecret"`
}

func (c *TestConfigWithUnexported) SecretFields() []string {
	return []string{"unexportedField", "NonStringField", "EmptySecret"}
}

type MySecretString string

func (s *MySecretString) SecretFields() []string {
	return []string{}
}

type MyConfig struct {
	SecretField string
}

func (c *MyConfig) SecretFields() []string {
	return []string{"SecretField"}
}

func TestLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	configFilePath := filepath.Join(tempDir, "TestConfig.yaml")
	fileContent := `
name: "test-name"
value: 123
enabled: true
nested:
  field1: "secret://nested-secret"
  field2: 3.14
deepNested:
  level2:
    deepField: "secret://deep-nested-secret"
SecretField: "secret://my-secret"
`
	err = os.WriteFile(configFilePath, []byte(fileContent), 0644)
	assert.NoError(t, err)

	_ = os.Setenv("PROPERTY_TESTCONFIG_NAME", "env-name")
	_ = os.Setenv("PROPERTY_TESTCONFIG_VALUE", "456")

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"cmd", "--enabled=false", "--nested-field1=cli-nested"}

	originalGetSecret := getSecretFunc
	defer func() { getSecretFunc = originalGetSecret }()
	getSecretFunc = func(name string) string {
		if name == "my-secret" {
			return "decrypted-secret"
		}
		if name == "nested-secret" {
			return "decrypted-nested-secret"
		}
		if name == "deep-nested-secret" {
			return "decrypted-deep-nested-secret"
		}
		return ""
	}

	config := &TestConfig{}

	originalGetFile := getFileFunc
	defer func() { getFileFunc = originalGetFile }()
	getFileFunc = func(appName string, cfg any) string {
		return configFilePath
	}

	Load("test-app", config)

	assert.Equal(t, "env-name", config.Name)            // Env should override file
	assert.Equal(t, 456, config.Value)                  // Env should override file
	assert.Equal(t, false, config.Enabled)              // Cli should override file
	assert.Equal(t, "cli-nested", config.Nested.Field1) // Cli should override file
	assert.Equal(t, 3.14, config.Nested.Field2)
	assert.Equal(t, "decrypted-secret", config.SecretField)
	assert.Equal(t, "decrypted-deep-nested-secret", config.DeepNested.Level2.DeepField)
}

func TestLoadFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	configFilePath := filepath.Join(tempDir, "TestConfig.yaml")
	fileContent := "name: \"test-name\"\nvalue: 123"
	err = os.WriteFile(configFilePath, []byte(fileContent), 0644)
	assert.NoError(t, err)

	config := &TestConfig{}
	err = LoadFileWithPath(config, configFilePath)
	assert.NoError(t, err)

	assert.Equal(t, "test-name", config.Name)
	assert.Equal(t, 123, config.Value)
}

func TestLoadFile_NotExist(t *testing.T) {
	config := &TestConfig{}
	err := LoadFileWithPath(config, "non-existent-file.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "", config.Name)
	assert.Equal(t, 0, config.Value)
}

func TestLoadFile_IsDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	config := &TestConfig{}
	err = LoadFileWithPath(config, tempDir)
	assert.NoError(t, err)

	assert.Equal(t, "", config.Name)
	assert.Equal(t, 0, config.Value)
}

func TestGetFile(t *testing.T) {
	path := getFile("my-app", &TestConfig{})
	assert.Contains(t, path, "my-app")
	assert.Contains(t, path, "TestConfig.yaml")
}

func TestGetFile_EmptyApplicationName(t *testing.T) {
	output := &bytes.Buffer{}
	log.Logger = log.Output(output)
	defer func() {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}()

	path := getFile("", &TestConfig{})
	assert.Contains(t, path, ".yaml") // Should still return a path, just without applicationName
	assert.Contains(t, output.String(), "application name is empty")
}

func TestLoadEnv(t *testing.T) {
	config := &TestConfig{}
	_ = os.Setenv("PROPERTY_TESTCONFIG_NAME", "env-name")
	_ = os.Setenv("PROPERTY_TESTCONFIG_VALUE", "456")
	_ = os.Setenv("PROPERTY_TESTCONFIG_ENABLED", "true")
	_ = os.Setenv("PROPERTY_TESTCONFIG_NESTED_FIELD1", "nested-env")
	_ = os.Setenv("PROPERTY_TESTCONFIG_NESTED_FIELD2", "1.23")

	LoadEnv(config)

	assert.Equal(t, "env-name", config.Name)
	assert.Equal(t, 456, config.Value)
	assert.Equal(t, true, config.Enabled)
	assert.Equal(t, "nested-env", config.Nested.Field1)
	assert.Equal(t, 1.23, config.Nested.Field2)
}

func TestSanitizeEnvKey(t *testing.T) {
	assert.Equal(t, "MYAPP_CONFIG", sanitizeEnvKey("MyApp.Config"))
	assert.Equal(t, "ANOTHER_CONFIG", sanitizeEnvKey("Another-Config"))
}

func TestLoadCli(t *testing.T) {
	config := &TestConfig{}
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd", "--name=cli-name", "--value=789", "--enabled=true"}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = gflag.ParseTo(config, fs)
	_ = fs.Parse(os.Args[1:])

	assert.Equal(t, "cli-name", config.Name)
	assert.Equal(t, 789, config.Value)
	assert.Equal(t, true, config.Enabled)
}

func TestLoadSecrets(t *testing.T) {
	config := &TestConfig{
		SecretField: "secret://my-secret",
		Nested: struct {
			Field1 string  `yaml:"field1"`
			Field2 float64 `yaml:"field2"`
		}{
			Field1: "secret://nested-secret",
		},
		DeepNested: struct {
			Level2 struct {
				DeepField string `yaml:"deepField"`
			} `yaml:"level2"`
		}{
			Level2: struct {
				DeepField string `yaml:"deepField"`
			}{
				DeepField: "secret://deep-nested-secret",
			},
		},
	}

	originalGetSecret := getSecretFunc
	defer func() { getSecretFunc = originalGetSecret }()
	getSecretFunc = func(name string) string {
		if name == "my-secret" {
			return "decrypted-value"
		}
		if name == "nested-secret" {
			return "decrypted-nested-value"
		}
		if name == "deep-nested-secret" {
			return "decrypted-deep-nested-value"
		}
		return ""
	}

	LoadSecrets(config)

	assert.Equal(t, "decrypted-value", config.SecretField)
	assert.Equal(t, "decrypted-nested-value", config.Nested.Field1)
	assert.Equal(t, "decrypted-deep-nested-value", config.DeepNested.Level2.DeepField)
}

func TestLoadSecrets_NoSecret(t *testing.T) {
	config := &TestConfig{
		SecretField: "not-a-secret",
	}
	LoadSecrets(config)
	assert.Equal(t, "not-a-secret", config.SecretField)
}

func TestLoadSecrets_NoSecretInterface(t *testing.T) {
	config := &struct{}{}
	LoadSecrets(config)
}

func TestLoadSecrets_NotPointer(t *testing.T) {
	config := TestConfig{}
	LoadSecrets(config)
}

func TestLoadSecrets_NotStruct(t *testing.T) {
	var config string
	LoadSecrets(&config)
}

func TestLoadSecrets_InvalidField(t *testing.T) {
	config := &TestConfig{}
	LoadSecrets(config)
}

func TestLoadSecrets_UnexportedField(t *testing.T) {
	config := &TestConfigWithUnexported{}
	LoadSecrets(config)
}

func TestLoadSecrets_NonStringField(t *testing.T) {
	config := &TestConfigWithUnexported{}
	LoadSecrets(config)
}

func TestLoadSecrets_EmptySecretName(t *testing.T) {
	config := &TestConfigWithUnexported{
		EmptySecret: "secret://",
	}
	LoadSecrets(config)
	assert.Equal(t, "secret://", config.EmptySecret)
}

func TestDecrypt(t *testing.T) {
	originalGetSecret := getSecretFunc
	defer func() { getSecretFunc = originalGetSecret }()
	getSecretFunc = func(name string) string {
		return "decrypted"
	}

	assert.Equal(t, "decrypted", Decrypt("any-secret"))
}

func TestLoadSecrets_PointerToNonStruct(t *testing.T) {
	var s MySecretString
	output := &bytes.Buffer{}
	log.Logger = log.Output(output)
	defer func() {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}()

	LoadSecrets(&s)
	assert.Contains(t, output.String(), "LoadSecrets expects a pointer to a struct; got: string")
}

func TestLoadSecrets_NilPointerToStruct(t *testing.T) {
	var nilConfig *MyConfig = nil

	output := &bytes.Buffer{}
	log.Logger = log.Output(output)
	defer func() {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}()

	LoadSecrets(nilConfig)
	assert.Contains(t, output.String(), "LoadSecrets expects a pointer to a struct; got: ptr")
}

type RequiredFieldsConfig struct {
	RequiredField string
	OptionalField string
}

func (c *RequiredFieldsConfig) RequiredFields() []string {
	return []string{"RequiredField"}
}

func TestValidateRequiredFields_ValidConfig(t *testing.T) {
	config := &RequiredFieldsConfig{
		RequiredField: "has-value",
		OptionalField: "",
	}
	validateRequiredFields(config)
}

func TestValidateRequiredFields_MissingField(t *testing.T) {
	config := &RequiredFieldsConfig{
		RequiredField: "",
		OptionalField: "value",
	}
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()
	validateRequiredFields(config)
}

func TestValidateRequiredFields_NonExistentField(t *testing.T) {
	config := &RequiredFieldsConfig{
		RequiredField: "value",
	}
	validateRequiredFields(config)
	assert.Equal(t, "value", config.RequiredField)
}

func TestValidateRequiredFields_NoImplementation(t *testing.T) {
	config := &TestConfig{}
	validateRequiredFields(config)
}

func TestValidateRequiredFields_NilPointer(t *testing.T) {
	var config *RequiredFieldsConfig = nil
	validateRequiredFields(config)
}

func TestValidateRequiredFields_NonStructPointer(t *testing.T) {
	type ConfigAsString string
	var config ConfigAsString = "test"
	validateRequiredFields(&config)
}

func TestValidateRequiredFields_MultipleMissingFields(t *testing.T) {
	output := &bytes.Buffer{}
	log.Logger = log.Output(output)
	defer func() {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}()

	type MultiFieldConfig struct {
		Field1 string
		Field2 string
		Field3 string
	}

	cfg := &MultiFieldConfig{}
	validateRequiredFields(cfg)
}

type ConfigForFieldValidation struct {
	StringField     string
	IntField        int
	BoolField       bool
	UnexportedField string
	EmptyString     string
}

func (c *ConfigForFieldValidation) RequiredFields() []string {
	return []string{"StringField", "IntField", "BoolField", "UnexportedField", "NonExistent"}
}

func TestIsValid_WithValidStringField(t *testing.T) {
	config := &ConfigForFieldValidation{
		StringField: "valid-value",
	}
	result := isValid(config, reflect.ValueOf(config).Elem(), "StringField")
	assert.True(t, result)
}

func TestIsValid_WithEmptyStringField(t *testing.T) {
	config := &ConfigForFieldValidation{
		StringField: "",
	}
	result := isValid(config, reflect.ValueOf(config).Elem(), "StringField")
	assert.False(t, result)
}

func TestIsValid_WithNonStringField(t *testing.T) {
	config := &ConfigForFieldValidation{}
	result := isValid(config, reflect.ValueOf(config).Elem(), "IntField")
	assert.False(t, result)
}

func TestIsValid_WithNonExistentField(t *testing.T) {
	config := &ConfigForFieldValidation{
		StringField: "value",
	}
	result := isValid(config, reflect.ValueOf(config).Elem(), "NonExistent")
	assert.False(t, result)
}

func TestLoadEnv_PartialEnvVars(t *testing.T) {
	config := &TestConfig{}
	_ = os.Setenv("PROPERTY_TESTCONFIG_NAME", "test-value")
	_ = os.Unsetenv("PROPERTY_TESTCONFIG_VALUE")

	LoadEnv(config)
	assert.Equal(t, "test-value", config.Name)
}

func TestLoadEnv_InvalidEnvVarType(t *testing.T) {
	config := &TestConfig{}
	_ = os.Setenv("PROPERTY_TESTCONFIG_VALUE", "not-a-number")

	LoadEnv(config)
}

type DeepConfig struct {
	Level1 struct {
		Level2 struct {
			Level3 struct {
				SecretField string
			}
		}
	}
}

func (c *DeepConfig) SecretFields() []string {
	return []string{"Level1.Level2.Level3.SecretField"}
}

func TestLoadSecrets_DeepNesting(t *testing.T) {
	config := &DeepConfig{}
	config.Level1.Level2.Level3.SecretField = "secret://deep-secret"

	originalGetSecret := getSecretFunc
	defer func() { getSecretFunc = originalGetSecret }()
	getSecretFunc = func(name string) string {
		if name == "deep-secret" {
			return "decrypted-deep-value"
		}
		return ""
	}

	LoadSecrets(config)
	assert.Equal(t, "decrypted-deep-value", config.Level1.Level2.Level3.SecretField)
}

func TestGetField_NestedPath(t *testing.T) {
	config := &TestConfig{
		Nested: struct {
			Field1 string  `yaml:"field1"`
			Field2 float64 `yaml:"field2"`
		}{
			Field1: "nested-value",
			Field2: 3.14,
		},
	}

	val := reflect.ValueOf(config).Elem()
	field := getField(val, "Nested.Field1")
	assert.True(t, field.IsValid())
	assert.Equal(t, "nested-value", field.String())
}

func TestGetField_SingleLevelPath(t *testing.T) {
	config := &TestConfig{
		Name: "test-value",
	}

	val := reflect.ValueOf(config).Elem()
	field := getField(val, "Name")
	assert.True(t, field.IsValid())
	assert.Equal(t, "test-value", field.String())
}

type CliConfig struct {
	Name  string
	Value int
}

func TestLoadCli_NoArgs(t *testing.T) {
	config := &CliConfig{}
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd"}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = gflag.ParseTo(config, fs)
	_ = fs.Parse(os.Args[1:])
}

func TestSanitizeEnvKey_SpecialCharacters(t *testing.T) {
	assert.Equal(t, "MYAPP_CONFIG", sanitizeEnvKey("MyApp.Config"))
	assert.Equal(t, "ANOTHER_CONFIG", sanitizeEnvKey("Another-Config"))
	assert.Equal(t, "APP_WITH_UNDERSCORES", sanitizeEnvKey("App_With_Underscores"))
	assert.Equal(t, "NUMBERS123", sanitizeEnvKey("Numbers123"))
	assert.Equal(t, "MULTIPLE_SEPS", sanitizeEnvKey("Multiple...___SEPS"))
	assert.Equal(t, "ONLY_ALPHANUMERIC_", sanitizeEnvKey("ONLY-ALPHANUMERIC!@#"))
}
