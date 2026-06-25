package main

import (
	"regexp"
)

func isValidCommandName(name string) bool {
	if len(name) == 0 || len(name) > 256 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-.]+$`, name)
	return matched
}
