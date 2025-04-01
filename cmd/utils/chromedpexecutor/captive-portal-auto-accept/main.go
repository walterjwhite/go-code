package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/headless"

	"github.com/chromedp/chromedp"
)

const (
	connectivityCheckUrl = "http://lxer.com"
)

func init() {
	application.Configure()
}

func main() {
	s := headless.New(application.Context)

	session.Execute(s, chromedp.Navigate(connectivityCheckUrl))
	chromedpexecutor.Screenshot(s.Context(), "/tmp/0.connectivity-check-fetch.png")

	acceptTerms(s)
	chromedpexecutor.Screenshot(s.Context(), "/tmp/1.connectivity-check-accept-terms.png")

	continueToInternet(s)
	chromedpexecutor.Screenshot(s.Context(), "/tmp/2.connectivity-check-continue.png")
}

func acceptTerms(s *headless.HeadlessChromeDPSession) {
	findAndClickElementByType(s, "input[type=\"checkbox\"]")
}

func continueToInternet(s *headless.HeadlessChromeDPSession) {
	findAndClickElementByType(s, "button")
}

func findAndClickElementByType(s *headless.HeadlessChromeDPSession, elementType string) {
	var elements []string
	session.Execute(s, chromedp.Evaluate(fmt.Sprintf("Array.from(document.querySelectorAll('%s')).map(element => element.outerHTML)", elementType), &elements))

	if len(elements) == 1 {
		log.Info().Msg("Found exactly 1 element, clicking it")
		elementIndex := 0

		session.Execute(s,
			chromedp.Click(fmt.Sprintf(`input[type="%s"]:nth-of-type(%d)`, elementType, elementIndex+1)),
		)
	} else {
		log.Warn().Msgf("Found %v elements, but expected just 1", len(elements))

		for _, element := range elements {
			log.Info().Msgf("Found element: %v", element)
		}
	}
}
