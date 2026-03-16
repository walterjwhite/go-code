package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidCharacter(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase z", 'z', true},
		{"lowercase m", 'm', true},

		{"uppercase A", 'A', true},
		{"uppercase Z", 'Z', true},
		{"uppercase M", 'M', true},

		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"digit 5", '5', true},

		{"hyphen", '-', true},
		{"underscore", '_', true},
		{"dot", '.', true},
		{"slash", '/', true},
		{"colon", ':', true},
		{"comma", ',', true},
		{"space", ' ', true},

		{"exclamation", '!', false},
		{"at sign", '@', false},
		{"hash", '#', false},
		{"dollar", '$', false},
		{"percent", '%', false},
		{"caret", '^', false},
		{"ampersand", '&', false},
		{"asterisk", '*', false},
		{"parenthesis open", '(', false},
		{"parenthesis close", ')', false},
		{"pipe", '|', false},
		{"backslash", '\\', false},
		{"semicolon", ';', false},
		{"quote", '"', false},
		{"single quote", '\'', false},
		{"backtick", '`', false},
		{"less than", '<', false},
		{"greater than", '>', false},
		{"equals", '=', false},
		{"plus", '+', false},
		{"question mark", '?', false},
		{"newline", '\n', false},
		{"tab", '\t', false},
		{"carriage return", '\r', false},
		{"null", '\x00', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCharacter(tt.input)
			if result != tt.expected {
				t.Errorf("isValidCharacter(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateArgument(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{"simple word", "hello", false, ""},
		{"mixed case", "HelloWorld", false, ""},
		{"with numbers", "test123", false, ""},
		{"with hyphen", "test-case", false, ""},
		{"with underscore", "test_case", false, ""},
		{"with dot", "test.case", false, ""},
		{"with slash", "test/case", false, ""},
		{"with colon", "test:case", false, ""},
		{"with comma", "test,case", false, ""},
		{"with space", "test case", false, ""},
		{"complex valid", "git,push origin", false, ""},
		{"alphanumeric mix", "abc123XYZ", false, ""},
		{"path like", "path/to/file.txt", false, ""},
		{"url like", "http://example.com:8080", false, ""},

		{"empty string", "", true, "argument cannot be empty"},

		{"too long", string(make([]byte, 1001)), true, "argument too long"},

		{"with exclamation", "test!", true, "invalid character"},
		{"with at sign", "test@case", true, "invalid character"},
		{"with hash", "test#case", true, "invalid character"},
		{"with dollar", "test$case", true, "invalid character"},
		{"with pipe", "test|case", true, "invalid character"},
		{"with semicolon", "test;case", true, "invalid character"},
		{"with newline", "test\ncase", true, "invalid character"},
		{"with tab", "test\tcase", true, "invalid character"},
		{"with null byte", "test\x00case", true, "invalid character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArgument(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateArgument(%q) expected error, got nil", tt.input)
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateArgument(%q) error = %q, expected to contain %q", tt.input, err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateArgument(%q) unexpected error: %v", tt.input, err)
				}
			}
		})
	}
}

func TestParseAndValidateArgs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectedLen int
	}{
		{"empty string", "", false, 0},
		{"single arg", "hello", false, 1},
		{"two args", "hello,world", false, 2},
		{"three args", "one,two,three", false, 3},
		{"with spaces", "hello, world , test", false, 3},
		{"complex args", "git,push,origin,master", false, 4},

		{"too many args", createManyArgs(101), true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				t.Skip("Skipping test that would call log.Fatalf")
				return
			}

			result := parseAndValidateArgs(tt.input)
			if len(result) != tt.expectedLen {
				t.Errorf("parseAndValidateArgs(%q) returned %d args, expected %d", tt.input, len(result), tt.expectedLen)
			}
		})
	}
}

func createManyArgs(count int) string {
	args := make([]string, count)
	for i := range count {
		args[i] = "arg"
	}
	return join(args, ",")
}

func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	var result strings.Builder
	result.WriteString(strs[0])
	for i := 1; i < len(strs); i++ {
		result.WriteString(sep + strs[i])
	}
	return result.String()
}

func TestValidateFilePath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"valid file in cwd", "test.txt", false},
		{"valid file with relative path", "./test.txt", false},
		{"valid file with absolute path in cwd", filepath.Join(cwd, "test.txt"), false},

		{"directory traversal", "../test.txt", true},
		{"absolute path outside cwd", "/tmp/test.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				t.Skip("Skipping test that would call log.Fatalf")
				return
			}

			tmpFile := filepath.Join(cwd, "test_validate_filepath.txt")
			if tt.input == "test.txt" || tt.input == "./test.txt" {
				tmpFile = filepath.Join(cwd, "test.txt")
			} else if tt.input == filepath.Join(cwd, "test.txt") {
				tmpFile = tt.input
			}

			f, err := os.Create(tmpFile)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			if err := f.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}
			defer func() { _ = os.Remove(tmpFile) }()

			result := validateFilePath(tt.input)
			if result == "" {
				t.Errorf("validateFilePath(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tests := []struct {
		name        string
		fileSize    int64
		expectError bool
	}{
		{"small file", 1024, false},                        // 1KB
		{"medium file", 1024 * 1024, false},                // 1MB
		{"large file under limit", 9 * 1024 * 1024, false}, // 9MB
		{"file at limit", 10 * 1024 * 1024, false},         // 10MB (exact limit)
		{"file over limit", 11 * 1024 * 1024, true},        // 11MB (over limit)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(cwd, "test_filesize.tmp")
			f, err := os.Create(tmpFile)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			if tt.fileSize > 0 {
				data := make([]byte, tt.fileSize)
				if _, err := f.Write(data); err != nil {
					_ = f.Close()
					_ = os.Remove(tmpFile)
					t.Fatalf("Failed to write to temp file: %v", err)
				}
			}
			if err := f.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}
			defer func() { _ = os.Remove(tmpFile) }()

			if tt.expectError {
				t.Skip("Skipping test that would call log.Fatalf")
				return
			}

			validateFileSize(tmpFile)
		})
	}
}

func TestServiceFactory(t *testing.T) {
	t.Run("NewServiceFactory", func(t *testing.T) {
		config := &ServiceConfig{
			ProjectID:     "test-project",
			UseEncryption: true,
			EncryptionKey: "test-key",
			Topic:         "test-topic",
			Subscription:  "test-sub",
			RegistryFile:  "test.json",
		}

		factory := NewServiceFactory(config)
		if factory == nil {
			t.Fatal("NewServiceFactory returned nil")
		}
		if factory.config != config {
			t.Error("ServiceFactory config not set correctly")
		}
	})

	t.Run("CreateEventService", func(t *testing.T) {
		t.Skip("Requires mocking of pubsub.NewClient")
	})
}

func TestConstants(t *testing.T) {
	t.Run("maxFileSize", func(t *testing.T) {
		expected := 10 * 1024 * 1024 // 10MB
		if maxFileSize != expected {
			t.Errorf("maxFileSize = %d, expected %d", maxFileSize, expected)
		}
	})

	t.Run("maxEvents", func(t *testing.T) {
		expected := 10000
		if maxEvents != expected {
			t.Errorf("maxEvents = %d, expected %d", maxEvents, expected)
		}
	})

	t.Run("maxArgs", func(t *testing.T) {
		expected := 100
		if maxArgs != expected {
			t.Errorf("maxArgs = %d, expected %d", maxArgs, expected)
		}
	})

	t.Run("maxArgLength", func(t *testing.T) {
		expected := 1000
		if maxArgLength != expected {
			t.Errorf("maxArgLength = %d, expected %d", maxArgLength, expected)
		}
	})

	t.Run("command constants", func(t *testing.T) {
		if commandList != "list" {
			t.Errorf("commandList = %q, expected %q", commandList, "list")
		}
		if commandRespond != "respond" {
			t.Errorf("commandRespond = %q, expected %q", commandRespond, "respond")
		}
		if commandListen != "listen" {
			t.Errorf("commandListen = %q, expected %q", commandListen, "listen")
		}
	})
}

func TestRegisterEvents(t *testing.T) {
	t.Skip("Requires mocking of service.EventService")
}

func TestLoadEventsFromFile(t *testing.T) {
	t.Skip("Requires mocking of service.EventService and file system")
}

func TestEdgeCases(t *testing.T) {
	t.Run("boundary length argument", func(t *testing.T) {
		exactMax := string(make([]byte, maxArgLength))
		for i := range exactMax {
			exactMax = exactMax[:i] + "a" + exactMax[i+1:]
		}
		err := validateArgument(exactMax)
		if err != nil {
			t.Errorf("validateArgument with exactly %d chars failed: %v", maxArgLength, err)
		}
	})

	t.Run("one over boundary length argument", func(t *testing.T) {
		overMax := string(make([]byte, maxArgLength+1))
		for i := range overMax {
			overMax = overMax[:i] + "a" + overMax[i+1:]
		}
		err := validateArgument(overMax)
		if err == nil {
			t.Errorf("validateArgument with %d chars should have failed", maxArgLength+1)
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		unicodeTests := []struct {
			name        string
			input       string
			expectError bool
		}{
			{"emoji", "test😀", true},
			{"chinese", "test 中文", true},
			{"arabic", "test عربي", true},
			{"cyrillic", "test русский", true},
		}

		for _, tt := range unicodeTests {
			t.Run(tt.name, func(t *testing.T) {
				err := validateArgument(tt.input)
				if tt.expectError && err == nil {
					t.Errorf("validateArgument(%q) should have failed for unicode", tt.input)
				}
			})
		}
	})

	t.Run("whitespace variations", func(t *testing.T) {
		whitespaceTests := []struct {
			name        string
			input       string
			expectError bool
		}{
			{"leading space", " test", false},
			{"trailing space", "test ", false},
			{"multiple spaces", "test  test", false},
			{"tab character", "test\ttest", true},
			{"newline", "test\ntest", true},
			{"carriage return", "test\rtest", true},
		}

		for _, tt := range whitespaceTests {
			t.Run(tt.name, func(t *testing.T) {
				err := validateArgument(tt.input)
				if tt.expectError && err == nil {
					t.Errorf("validateArgument(%q) should have failed", tt.input)
				}
				if !tt.expectError && err != nil {
					t.Errorf("validateArgument(%q) unexpected error: %v", tt.input, err)
				}
			})
		}
	})
}

func TestSecurityScenarios(t *testing.T) {
	t.Run("injection attempts", func(t *testing.T) {
		injectionTests := []struct {
			name  string
			input string
		}{
			{"sql injection", "test'; DROP TABLE users;--"},
			{"command injection", "test; rm -rf /"},
			{"path traversal", "../../../etc/passwd"},
			{"script injection", "<script>alert('xss')</script>"},
			{"null byte injection", "test\x00injection"},
		}

		for _, tt := range injectionTests {
			t.Run(tt.name, func(t *testing.T) {
				err := validateArgument(tt.input)
				if err == nil {
					t.Errorf("validateArgument(%q) should have rejected potentially dangerous input", tt.input)
				}
			})
		}
	})
}
