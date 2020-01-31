package main

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/email"
)

var (
	//serverAddressFlag = flag.String("ServerAddress", "imap.gmail.com:993", "Server Address")
	//logoutTimeoutFlag = flag.Int("LogoutTimeout", 30, "Logout Timeout in seconds")
	usernameFlag = flag.String("Username", "", "Username")
	passwordFlag = flag.String("Password", "", "Password")
	folderFlag   = flag.String("Folder", "INBOX", "Folder Name")
)

func init() {
	application.Configure()
}

func main() {
	emailSenderAccount := &email.EmailSenderAccount{Username: *usernameFlag,
		Password:   *passwordFlag,
		ImapServer: &email.EmailServer{Host: "imap.gmail.com", Port: 993}}

	emailSession := emailSenderAccount.Connect()
	emailSession.ReadAsync(*folderFlag)
}
