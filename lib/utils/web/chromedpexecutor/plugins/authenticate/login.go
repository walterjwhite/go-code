package authenticate

import (
	"context"

	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/property"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

func (s *Session) Login(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	//s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(chromedp.Navigate(*s.Website.Url))

	log.Debug().Msg("fetched")
	for _, fieldGroup := range s.Website.FieldGroups {
		for _, field := range fieldGroup.Fields {
			log.Debug().Msgf("get value: %s", *field.Id)

			value, err := getValue(s.Credentials, field.Id)
			logging.Panic(err)

			log.Debug().Msgf("executing: %s / %s", *field.Selector, value)
			// s.chromedpsession.Execute(chromedp.SendKeys(*field.Selector, value, chromedp.ByID))
			s.chromedpsession.Execute(chromedp.SendKeys(*field.Selector, value))
			log.Debug().Msgf("executed: %s", *field.Selector)
		}

		log.Debug().Msg("submitting")
		selector := getSubmitSelector(fieldGroup)
		s.chromedpsession.Execute(chromedp.Submit(*selector))
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
