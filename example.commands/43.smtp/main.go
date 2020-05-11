package main

import (
	"flag"
	"github.com/emersion/go-message/mail"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/email"
)

var (
	//serverAddressFlag = flag.String("ServerAddress", "imap.gmail.com:993", "Server Address")
	//logoutTimeoutFlag = flag.Int("LogoutTimeout", 30, "Logout Timeout in seconds")
	usernameFlag = flag.String("Username", "", "Username")
	passwordFlag = flag.String("Password", "", "Password")
)

func init() {
	application.Configure()
}

func main() {
	emailSenderAccount := &email.EmailSenderAccount{Username: *usernameFlag,
		Password:     *passwordFlag,
		EmailAddress: &mail.Address{Address: *usernameFlag},
		SmtpServer:   &email.EmailServer{Host: "smtp.gmail.com", Port: 465}}

	// TODO: automatically set from based on the sender account ...
	emailMessage := &email.EmailMessage{
		To: []*mail.Address{
			{Address: *usernameFlag}},
		Subject: "Testing",
		Body:    "Testing - body"}
	emailSenderAccount.Send(emailMessage)
}
