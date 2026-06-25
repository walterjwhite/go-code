package graphical

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/until"
	"time"
)

func (i *ImageMatch) WaitUntilMatched(pctx context.Context, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(pctx, timeout)
	defer cancel()

	err := until.Until(ctx, 250*time.Millisecond, i.Matches)
	if err != nil {
		log.Warn().Msgf(".isApplicationIconPresent.error - %v", err)
		return false, err
	}

	return i.Match != nil, nil
}
