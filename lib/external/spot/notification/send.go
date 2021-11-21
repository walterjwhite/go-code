package notification

import (
	"github.com/walterjwhite/go-code/lib/net/email"
)

const (
	dateFormat = "2006/01/02 15:04:05 -0700 MST"
)

func (c *Notification) Send() {
	c.EmailMessage.Attachments = []*email.EmailAttachment{c.getTrackAsAttachment()}
	c.addAttachments()

	c.prepareTemplateContext()
	c.doTemplate()

	c.EmailSenderAccount.Send(c.EmailMessage)
}
