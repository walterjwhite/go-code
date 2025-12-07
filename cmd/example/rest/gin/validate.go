package main

import (
	"net/mail"
)

func validateEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
