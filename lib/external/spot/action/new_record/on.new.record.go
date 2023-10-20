package new_record

import (
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/external/spot/notification"
	"path/filepath"
)

func (c *Configuration) OnNewRecord(old, new *data.Record) {
	templateName := filepath.Join("new-record", string(new.MessageType))

	if notification.Exists(c.Session, templateName) {
		notification.New(c.Session, new, templateName).Send()
	}
}
