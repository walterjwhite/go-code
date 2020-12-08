package notification

import (
	"bytes"

	"github.com/walterjwhite/go/lib/application/logging"

	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go/lib/net/email"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Notification) addAttachments() {
	for _, filename := range c.Filenames {
		f := c.getAttachmentFilename(filename)

		data, err := ioutil.ReadFile(f)
		logging.Panic(err)

		fileAttachment := &email.EmailAttachment{Name: "*." + filepath.Base(filename), Data: bytes.NewBuffer(data)}
		c.EmailMessage.Attachments = append(c.EmailMessage.Attachments, fileAttachment)
	}
}

func (c *Notification) getAttachmentFilename(filename string) string {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return filename
	}

	// not absolute, try relative
	relativeFilename := filepath.Join(c.Session.SessionPath, ".notifications", filename)
	_, err = os.Stat(relativeFilename)
	if !os.IsNotExist(err) {
		return relativeFilename
	}

	expandedFilename, err := homedir.Expand(filename)
	logging.Panic(err)

	return expandedFilename
}
