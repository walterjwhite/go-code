package notification

import (
	"github.com/walterjwhite/go-code/lib/application/property"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/io/yaml"
	"github.com/walterjwhite/go-code/lib/net/email"

	"os"
	"path/filepath"
)

// read from file (in session path)
type Notification struct {
	TemplateName string

	Session *data.Session
	Record  *data.Record

	// email configuration
	EmailSenderAccount *email.EmailSenderAccount

	EmailMessage *email.EmailMessage

	// for simplicity, only support strings
	Context map[string]interface{}

	//ReferenceData *ReferenceData
	Filenames []string

	// TODO: implement interface (email, SMS, IM, etc.)
	//Notifiers []*Notifier
}

// TODO: make this generic, not tied to email
func New(s *data.Session, r *data.Record, templateName string) *Notification {
	n := &Notification{}
	yaml.Read(getTemplateName(s, templateName), n)

	property.LoadSecrets(n.EmailSenderAccount)

	n.Context = make(map[string]interface{})
	n.Session = s
	n.Record = r
	n.TemplateName = templateName

	n.loadReferences()

	return n
}

func getTemplateName(s *data.Session, templateName string) string {
	return filepath.Join(property.GetConfigurationDirectory("spot", s.FeedId, "notifications"), templateName+".yaml")
}

func Exists(s *data.Session, templateName string) bool {
	filename := getTemplateName(s, templateName)

	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
