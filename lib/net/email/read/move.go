package write

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-move"
)

func (s *EmailSession) Move(msg *imap.Message, destination string) error {
	if err := validateFolderName(destination); err != nil {
		return fmt.Errorf("invalid destination folder name: %w", err)
	}

	mc := move.NewClient(s.client)

	set := new(imap.SeqSet)
	set.AddNum(msg.Uid)

	return mc.UidMove(set, destination)
}
