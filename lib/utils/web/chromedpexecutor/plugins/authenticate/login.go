package authenticate

import (
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/property"
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
			// TODO: use wrapper which introduces a delay, waits until visible, etc.
			logging.Panic(chromedp.Run(s.chromedpsession.Context(), chromedp.SendKeys(*field.Selector, value)))
			log.Debug().Msgf("executed: %s", *field.Selector)
		}

		log.Debug().Msg("submitting")
		selector := getSubmitSelector(fieldGroup)
		logging.Panic(chromedp.Run(s.chromedpsession.Context(), chromedp.Submit(*selector)))
		log.Debug().Msg("submitted")
	}
}

func getSubmitSelector(fieldGroup *FieldGroup) *string {
	if fieldGroup.SubmitSelector != nil {
		return fieldGroup.SubmitSelector
	}

	return fieldGroup.Fields[len(fieldGroup.Fields)-1].Selector
}

func getValue(c *Credentials, FieldId *string) (string, error) {
	for _, secret := range c.Secrets {
		if strings.Compare(*secret.FieldId, *FieldId) == 0 {
			return property.Decrypt(*secret.SecretKey), nil
		}
	}

	return "", fmt.Errorf("Field not found: %s", *FieldId)
}
