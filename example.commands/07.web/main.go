package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"github.com/ddo/rq"
	"net/http"
)

func main() {
	r := rq.Get("https://pnc.com")

	// send with golang default HTTP client
	req, err := r.ParseRequest()
	logging.Panic(err)

	res, err := http.DefaultClient.Do(req)
	logging.Panic(err)

	defer res.Body.Close()

	log.Info().Msgf("Response:\n%v", res)
}
