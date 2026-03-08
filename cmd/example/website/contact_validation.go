package main

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	maxMessageLength = 5000
	maxNameLength    = 120
	maxEmailLength   = 254
	maxSubjectLength = 200
)

func validateContactRequest(contactRequest *ContactRequest) error {
	if contactRequest == nil {
		return errors.New("invalid request payload")
	}

	contactRequest.Name = strings.TrimSpace(contactRequest.Name)
	contactRequest.Email = strings.TrimSpace(contactRequest.Email)
	contactRequest.Subject = strings.TrimSpace(contactRequest.Subject)
	contactRequest.Message = strings.TrimSpace(contactRequest.Message)

	if contactRequest.Name == "" ||
		contactRequest.Email == "" ||
		contactRequest.Subject == "" ||
		contactRequest.Message == "" {
		return errors.New("all fields are required")
	}

	if utf8.RuneCountInString(contactRequest.Name) > maxNameLength {
		return errors.New("name exceeds maximum length")
	}

	if utf8.RuneCountInString(contactRequest.Email) > maxEmailLength {
		return errors.New("email exceeds maximum length")
	}

	if utf8.RuneCountInString(contactRequest.Subject) > maxSubjectLength {
		return errors.New("subject exceeds maximum length")
	}

	if utf8.RuneCountInString(contactRequest.Message) > maxMessageLength {
		return errors.New("message exceeds maximum length of 5000 characters")
	}

	if !validateEmailAddress(contactRequest.Email) {
		return errors.New("invalid email address")
	}

	if hasCRLF(contactRequest.Subject) || hasCRLF(contactRequest.Email) || hasCRLF(contactRequest.Name) {
		return errors.New("invalid characters in request")
	}

	return nil
}

func validateEmailAddress(email string) bool {
	parsed, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if parsed.Name != "" {
		return false
	}

	return parsed.Address == email
}

func hasCRLF(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}
