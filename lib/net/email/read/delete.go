package write

import (
	"github.com/emersion/go-imap"
)

func (s *EmailSession) Delete(msg *imap.Message) error {
	set := new(imap.SeqSet)
	set.AddNum(msg.Uid)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []any{imap.DeletedFlag}
	err := s.client.Store(set, item, flags, nil)
	if err != nil {
		return err
	}

	return s.client.Expunge(nil)
}
