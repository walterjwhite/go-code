package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/elasticsearch"
	"github.com/walterjwhite/go-application/libraries/email"
	"time"
)

func init() {
	application.Configure()
}

func main() {
	c := elasticsearch.NewDefaultClient()
	b := c.NewBatch(10, 5242880, 30*time.Second, 2)

	emailMessage := email.EmailMessage{To: []string{"<TO>"}, Subject: "Test Email", Body: "Test Body\n"}

	// integrate secrets here ...
	emailSenderAccount := &email.EmailSenderAccount{Username: "<USERNAME>", Password: "<PASSWORD>",
		EmailAddress: "<EMAILADDRESS>", Server: email.EmailServer{Host: "smtp.gmail.com", Port: 465}}

	b.Index( /*email.Id*/ "email.1", emailMessage)

	b.Flush()

	emailSenderAccount.Send(emailMessage)
}
