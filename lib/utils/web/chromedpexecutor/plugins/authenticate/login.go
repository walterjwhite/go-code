package authenticate

import (
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/rs/zerolog/log"

	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/property"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

func (s *Session) Login() {
	// handled outside of here
	// s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	//s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	logging.Panic(chromedp.Run(s.chromedpsession.Context(), chromedp.Navigate(*s.Website.Url)))

	log.Debug().Msgf("fetched: %v", s.Website)
	for _, fieldGroup := range s.Website.FieldGroups {
		for _, field := range fieldGroup.Fields {
			log.Debug().Msgf("get value: %s", *field.Id)

			value, err := getValue(s.Credentials, field.Id)
			logging.Panic(err)

			log.Debug().Msgf("executing: %s / %s", *field.Selector, value)
			// logging.Panic(chromedp.Run(s.chromedpsession.Context(), chromedp.Clear(*field.Selector)))
			// log.Debug().Msgf("cleared: %s", *field.Selector)

			logging.Panic(action.SendKeys(s.chromedpsession.Context(), *s.VisibleTimeout, s.LocateDelay, *field.Selector, value))

			log.Debug().Msgf("executed: %s", *field.Selector)
		}

		log.Debug().Msg("submitting")
		s.LocateDelay.Delay()

		s.submit(fieldGroup)

		log.Debug().Msg("submitted")
	}
}

func (s *Session) submit(fieldGroup *FieldGroup) {
	if fieldGroup.SubmitSelector != nil {
		logging.Panic(chromedp.Run(s.chromedpsession.Context(), chromedp.Submit(fieldGroup.SubmitSelector)))
		return
	}

	submitSelector := fieldGroup.Fields[len(fieldGroup.Fields)-1].Selector
	log.Debug().Msgf("submit selector: %v", *submitSelector)

	logging.Panic(action.SendKeys(s.chromedpsession.Context(), *s.VisibleTimeout, s.LocateDelay, *submitSelector, kb.Enter))
}

func getValue(c *Credentials, FieldId *string) (string, error) {
	for _, secret := range c.Secrets {
		if strings.Compare(*secret.FieldId, *FieldId) == 0 {
			return property.Decrypt(*secret.SecretKey), nil
		}
	}

	return "", fmt.Errorf("Field not found: %s", *FieldId)
}
