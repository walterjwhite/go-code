package email

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-move"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *EmailSession) Move(msg *imap.Message, destination string) {
	mc := move.NewClient(s.client)

	set := new(imap.SeqSet)
	set.AddNum(msg.Uid)

	logging.Panic(mc.UidMove(set, destination))
}
