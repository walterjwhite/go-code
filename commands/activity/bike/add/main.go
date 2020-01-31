package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/bike"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
	"strconv"
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 4 {
		log.Info().Msgf("Adding bike mileage: %v", flag.Args())

		d, err := strconv.ParseFloat(flag.Args()[2], 64)
		logging.Panic(err)

		index, err := strconv.Atoi(flag.Args()[1])
		logging.Panic(err)

		b := bike.Add(application.Context, flag.Args()[0], index, d, flag.Args()[3], flag.Args()[4:]...)
		log.Info().Msgf("added bike mileage: %v", b)
	} else {
		logging.Panic(fmt.Errorf("No bike mileage entered - expecting: YYYY/MM/DD index distance bike tags..."))
	}
}
