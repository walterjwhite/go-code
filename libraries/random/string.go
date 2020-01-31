package random

import (
	"errors"
	"github.com/walterjwhite/go-application/libraries/logging"
	"math/rand"
	"time"
)

const (
	defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func validateInput(length int, charset string) {
	if length <= 0 {
		logging.Panic(errors.New("Please enter a non-zero length"))
	}

	if len(charset) <= 0 {
		logging.Panic(errors.New("Please enter a non-empty charset"))
	}
}

func StringWithCharset(length int, charset string) string {
	validateInput(length, charset)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, defaultCharset)
}
