package main

import (
	"flag"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/price-checker"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/price-checker/amazon"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/price-checker/jensonusa"
)







var (
	url = flag.String("url", "", "url")
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	s := price_checker.New(application.Context, url)
	if strings.Contains(*url, "jensonusa.com") {
		s.PriceChecker = jensonusa.New(s.Session, url)
	} else if strings.Contains(*url, "amazon.com") {
		s.PriceChecker = amazon.New(s.Session, url)
	}

	price := s.Fetch()

	log.Info().Msgf("price: %s", price)
}
