package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/token/plugins/stdin"

	"github.com/walterjwhite/go-application/libraries/gmail"

	"flag"
)

var (
	// can only use a phone number so many times (10)
	phoneNumberFlag       = flag.String("PhoneNumber", "", "Phone Number")
	callInsteadOfTextFlag = flag.Bool("Call", false, "Call instead of text")
	phoneUsageFlag        = flag.Bool("PhoneUsage", false, "Do not use phone integration for this account")
)

func init() {
	application.Configure()
}

func main() {
	tokenProvider := &stdin.StdInReader{PromptMessage: "Please enter the verification code from Google.\n"}
	gmail.NewRandom(&gmail.PhonePreference{PhoneNumber: *phoneNumberFlag, Call: *callInsteadOfTextFlag, PhoneUsage: *phoneUsageFlag}).Create(application.Context, tokenProvider)
}
