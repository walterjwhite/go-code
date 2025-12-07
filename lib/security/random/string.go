package random

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func validateInput(length int, charset string) error {
	if length <= 0 {
		return errors.New("please enter a non-zero length")
	}

	if len(charset) <= 0 {
		return errors.New("please enter a non-empty charset")
	}

	return nil
}

func StringWithCharset(length int, charset string) (string, error) {
	if err := validateInput(length, charset); err != nil {
		return "", err
	}

	out := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = charset[n.Int64()]
	}

	return string(out), nil
}

func String(length int) (string, error) {
	return StringWithCharset(length, defaultCharset)
}
