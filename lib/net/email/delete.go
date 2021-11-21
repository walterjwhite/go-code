package email

import (
	"github.com/emersion/go-imap"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *EmailSession) Delete(msg *imap.Message) {
	//dc := move.NewClient(s.client)

	set := new(imap.SeqSet)
	set.AddNum(msg.Uid)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	logging.Panic(s.client.Store(set, item, flags, nil))

	logging.Panic(s.client.Expunge(nil))
}
