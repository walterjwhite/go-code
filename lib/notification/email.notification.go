package notification

import (
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/write"
)

type EmailNotification struct {
	EmailAccount *email.EmailAccount

	EmailMessage *email.EmailMessage
}

func (n *EmailNotification) Notify(notification Notification) error {
	n.EmailMessage.Subject = notification.Title
	n.EmailMessage.Body = notification.Description

	return write.Send(n.EmailAccount, n.EmailMessage)
}
