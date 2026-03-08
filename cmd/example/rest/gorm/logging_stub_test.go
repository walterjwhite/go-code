package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Nil(t *testing.T) {
	Error(nil)
}

func TestError_WithError(t *testing.T) {
	err := NewError("test error")
	Error(err)
	assert.Equal(t, "test error", err.Error())
}

func TestErrorWithNil(t *testing.T) {
	ErrorWithNil()
}

func TestErrorWithActual(t *testing.T) {
	err := NewError("actual error")
	ErrorWithActual(err)
	assert.Equal(t, "actual error", err.Error())
}

func TestNewError(t *testing.T) {
	err := NewError("new error")
	assert.NotNil(t, err)
	assert.Equal(t, "new error", err.Error())
}
