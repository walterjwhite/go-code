package main

import (
	"regexp"
	"unicode"
)

func isValidCommandName(name string) bool {
	if len(name) == 0 || len(name) > 256 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-.]+$`, name)
	return matched
}

func isValidArgument(arg string) bool {
	const maxArgLength = 4096
	if len(arg) > maxArgLength {
		return false
	}

	for _, r := range arg {
		if !isValidCharacter(r) {
			return false
		}
	}
	return true
}

func isValidCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) ||
		unicode.IsSpace(r) ||
		r == '.' || r == ',' || r == '-' || r == '_' ||
		r == '+' || r == '='
}
