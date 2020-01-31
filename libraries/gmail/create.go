package gmail

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/token"
)

func (a *Account) Create(ctx context.Context, tokenProvider token.TokenProvider) {
	a.validate()

	log.Info().Msgf("New GMAIL account: %v", *a)
	a.doCreate(ctx, tokenProvider)
}
