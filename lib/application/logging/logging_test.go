package logging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		Panic(nil)
	})

	assert.Panics(t, func() {
		Panic(errors.New("test error"))
	})

	assert.Panics(t, func() {
		Panic(errors.New("test error"), "context1", 123)
	})
}

func TestWarn(t *testing.T) {
	assert.NotPanics(t, func() {
		Warn(nil, "test message")
	})

	assert.NotPanics(t, func() {
		Warn(errors.New("test error"), "test message")
	})
}
