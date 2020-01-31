package notification

import (
	"github.com/walterjwhite/go-application/libraries/email"
)

type EmailNotification struct {
	EmailSenderAccount *email.EmailSenderAccount

	EmailMessage *email.EmailMessage
}

func (n *EmailNotification) Notify(notification Notification) {
	n.EmailMessage.Subject = notification.Title
	n.EmailMessage.Body = notification.Description

	n.EmailSenderAccount.Send(n.EmailMessage)
}
