package citrix

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	useLightVersionPromptId = "protocolhandler-welcome-useLightVersionLink"
)

func (s *Session) useLightVersion() {
	log.Info().Msgf("Session.useLightVersion - UseLightVersion: %v", s.UseLightVersion)

	if !s.UseLightVersion {
		return
	}

	select {
	case <-s.ctx.Done():
		log.Debug().Msg("Session.useLightVersion - context done")
	default:
	}

	if action.ExistsById(s.ctx, useLightVersionPromptId) {
		log.Info().Msg("Session.useLightVersion - switching to light version")
		logging.Warn(action.Execute(s.ctx,
			chromedp.Click(useLightVersionPromptId, chromedp.ByID),
		), "session.useLightVersion - error selecting use light version")
	}
}
