package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/part"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
	"strconv"
	"strings"
)

var (
	tagsFlag = flag.String("Tags", "", "Tags")
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 2 {
		log.Info().Msgf("Adding Part Provider: %v", flag.Args()[0])

		var tags []string
		if len(*tagsFlag) > 0 {
			tags = strings.Split(*tagsFlag, ",")
		}

		var price, tax, total float64
		var err error
		price, err = strconv.ParseFloat(flag.Args()[3], 32)
		logging.Panic(err)

		tax, err = strconv.ParseFloat(flag.Args()[4], 32)
		logging.Panic(err)

		total, err = strconv.ParseFloat(flag.Args()[5], 32)
		logging.Panic(err)

		q := part.AddProvider(application.Context, flag.Args()[0], flag.Args()[1], flag.Args()[2], float32(price), float32(tax), float32(total), tags)
		log.Info().Msgf("Added Part: %v", q)
	} else {
		logging.Panic(fmt.Errorf("No part added"))
	}
}
