package main

import (
	"net/url"
	"regexp"
	"strings"
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

	if looksLikeURL(arg) {
		return isValidYouTubeURL(arg)
	}

	for _, r := range arg {
		if !isValidCharacter(r) {
			return false
		}
	}
	return true
}

func looksLikeURL(arg string) bool {
	lower := strings.ToLower(arg)
	return strings.Contains(lower, "http://") || strings.Contains(lower, "https://")
}

func isValidYouTubeURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	if parsed.Scheme != "https" {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	allowedHosts := map[string]bool{
		"youtube.com":     true,
		"www.youtube.com": true,
		"youtu.be":        true,
		"m.youtube.com":   true,
	}
	if !allowedHosts[host] {
		return false
	}

	if parsed.Fragment != "" {
		return false
	}

	for key, values := range parsed.Query() {
		if !isValidURLComponent(key) {
			return false
		}
		for _, v := range values {
			if !isValidURLComponent(v) {
				return false
			}
		}
	}

	return true
}

func isValidURLComponent(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) &&
			r != '-' && r != '_' && r != '.' && r != '~' {
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
