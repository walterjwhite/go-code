package gmail

import (
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
)

type Gender int

const (
	firstNameMinimumLength = 2
	lastNameMinimumLength  = 2
	usernameMinimumLength  = 6
	usernameMaximumLength  = 30
	passwordMinimumLength  = 8
	phoneNumberLength      = 10
	minimumAge             = 10
	maximumAge             = 50

	deviationFactor = 2

	gmailBaseUrl = "https://accounts.google.com/SignUp?hl=en"
)

const (
	Male Gender = iota + 1
	Female
	RatherNotSay
	Custom
)

type PhonePreference struct {
	PhoneNumber string
	// call instead of text
	Call       bool
	PhoneUsage bool
}

type Account struct {
	FirstName string
	LastName  string
	Username  string
	Password  string

	BirthDate *BirthDate
	Gender    Gender

	PhonePreference *PhonePreference

	session *chromedpexecutor.ChromeDPSession
}

type BirthDate struct {
	Month int
	Day   int
	Year  int
}
