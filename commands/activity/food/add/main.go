package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/food"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
	"strconv"
	"time"
)

const (
	dateFormat = "2006/01/02"
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 3 {
		log.Info().Msgf("Adding food entry: %v", flag.Args())

		parsedDate, parseErr := time.Parse(dateFormat, flag.Args()[0])
		logging.Panic(parseErr)

		instance, err := strconv.Atoi(flag.Args()[1])
		logging.Panic(err)

		f := food.Add(parsedDate, instance, flag.Args()[2:]...)
		log.Info().Msgf("added food entry: %v", f)
	} else {
		logging.Panic(fmt.Errorf("No food entry entered"))
	}
}
