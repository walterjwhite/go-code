package logging

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"
)

func TestLogContextuals(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = log.Output(io.MultiWriter(&buf, os.Stderr))
	defer func() {
		log.Logger = originalLogger
	}()

	tests := []struct {
		name       string
		contextual []any
		expectLog  bool
	}{
		{
			name:       "nil contextuals",
			contextual: nil,
			expectLog:  false,
		},
		{
			name:       "empty contextuals",
			contextual: []any{},
			expectLog:  false,
		},
		{
			name:       "with contextuals",
			contextual: []any{"context1", 123},
			expectLog:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			logContextuals(tt.contextual...)
			if tt.expectLog {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestIsDevEnvironment(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		want   bool
	}{
		{
			name:   "development environment",
			envVal: "development",
			want:   true,
		},
		{
			name:   "production environment",
			envVal: "production",
			want:   false,
		},
		{
			name:   "empty environment",
			envVal: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENVIRONMENT", tt.envVal)

			assert.Equal(t, tt.want, isDevEnvironment())
		})
	}
}

func TestLogErrorMessage(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = log.Output(io.MultiWriter(&buf, os.Stderr))
	defer func() {
		log.Logger = originalLogger
	}()

	buf.Reset()
	err := errors.New("test error message")
	logErrorMessage(err)

	output := buf.String()
	assert.Contains(t, output, "level\":\"error")
	assert.Contains(t, output, "test error message")
}

func TestLogSecurityNote(t *testing.T) {
	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = log.Output(io.MultiWriter(&buf, os.Stderr))
	defer func() {
		log.Logger = originalLogger
	}()

	buf.Reset()
	logSecurityNote()

	output := buf.String()
	assert.Contains(t, output, "Stack trace unavailable in production for security reasons")
}

func TestError(t *testing.T) {
	type args struct {
		err         error
		contextuals []any
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "nil error, no panic",
			args:      args{err: nil},
			wantPanic: false,
		},
		{
			name:      "non-nil error, should not panic",
			args:      args{err: errors.New("test error")},
			wantPanic: false,
		},
		{
			name:      "non-nil error with contextuals, should not panic",
			args:      args{err: errors.New("test error with context"), contextuals: []any{"context1", 123}},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					Error(tt.args.err, tt.args.contextuals...)
				})
			} else {
				assert.NotPanics(t, func() {
					Error(tt.args.err, tt.args.contextuals...)
				})
			}
		})
	}
}

func TestWarn(t *testing.T) {


	var buf bytes.Buffer

	originalLogger := log.Logger // Save original logger

	log.Logger = log.Output(io.MultiWriter(&buf, os.Stderr)) // Redirect global logger output

	defer func() {

		log.Logger = originalLogger // Restore original logger

	}()

	tests := []struct {
		name string

		err error

		message string

		expectedLogs []string

		notExpectedLogs []string
	}{

		{

			name: "nil error, no log output",

			err: nil,

			message: "test message",
		},

		{

			name: "non-nil error, should log message without stack trace",

			err: errors.New("warning error"),

			message: "test message",

			expectedLogs: []string{

				"level\":\"warn",

				"message\":\"test message - warning error",
			},

			notExpectedLogs: []string{

				"level\":\"panic",

				"error\":\"warning error",

				"Stack trace:",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			buf.Reset() // Clear buffer for each test case

			Warn(tt.err, tt.message)

			output := buf.String()

			if tt.err == nil {

				assert.Empty(t, output, "Expected no log output for nil error")

			} else {

				for _, expected := range tt.expectedLogs {

					assert.Contains(t, output, expected)

				}

				for _, notExpected := range tt.notExpectedLogs {

					assert.NotContains(t, output, notExpected)

				}

			}

		})

	}

}
