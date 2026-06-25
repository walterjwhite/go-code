package random

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	defaultCharset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxStringLength = 1024 // Maximum allowed string length to prevent memory exhaustion
	minUniqueChars  = 2    // Minimum unique characters required for adequate entropy
)

func validateInput(length int, charset string) error {
	if length <= 0 {
		return errors.New("please enter a non-zero length")
	}

	if length > maxStringLength {
		return errors.New("length exceeds maximum allowed size")
	}

	if len(charset) < minUniqueChars {
		return errors.New("charset must contain at least 2 unique characters for adequate entropy")
	}

	charSet := make(map[rune]bool)
	for _, c := range charset {
		if charSet[c] {
			return errors.New("charset contains duplicate characters which reduces entropy")
		}
		charSet[c] = true
	}

	return nil
}

func StringWithCharset(length int, charset string) (string, error) {
	if err := validateInput(length, charset); err != nil {
		return "", err
	}

	out := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range length {
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
