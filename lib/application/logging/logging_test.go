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

func TestError(t *testing.T) {
	type args struct {
		err         error
		contextuals []interface{}
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
			name:      "non-nil error, should panic",
			args:      args{err: errors.New("test error")},
			wantPanic: true,
		},
		{
			name:      "non-nil error with contextuals, should panic",
			args:      args{err: errors.New("test error with context"), contextuals: []interface{}{"context1", 123}},
			wantPanic: true,
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

			name: "non-nil error, should log message and stack trace",

			err: errors.New("warning error"),

			message: "test message",

			expectedLogs: []string{

				"level\":\"warn",

				"message\":\"test message - warning error",

				"Stack trace:",
			},

			notExpectedLogs: []string{

				"level\":\"panic",

				"error\":\"warning error",
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
