package write

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-move"
)

func (s *EmailSession) Move(msg *imap.Message, destination string) error {
	mc := move.NewClient(s.client)

	set := new(imap.SeqSet)
	set.AddNum(msg.Uid)

	return mc.UidMove(set, destination)
}
