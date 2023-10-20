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

//*[@id="price"]/text()
////*[@id="js_attskuselector_BR001027"]/div[2]/div[2]/div/div/div/text()
//*[@id="txtAddToCart"]

// jensonusa:
// 	price: https://www.jensonusa.com/api/product/GetExtendOffers?productId=BR001048%20160
// 	inStock: https://www.jensonusa.com/jensonProduct/getStockDetails?productCode=BR001048&variationCode=BR001048%20160&productStatus=

// amazon:
// 	HTML parsing

// performancebike:
// 	HTML parsing

// universalcycles.com:
// 	HTML parsing


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

	// price := amazon.FetchAmazon(application.Context, url)
	price := s.Fetch()

	log.Info().Msgf("price: %s", price)
}
