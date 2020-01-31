package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/bike"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/olivere/elastic.v7"
	"time"

	"flag"
	"github.com/rs/zerolog/log"
)

const (
	dateFormat = "2006/01/02"
)

var (
	bikeFlag  = flag.String("Bike", "", "Bike Name")
	startFlag = flag.String("StartDate", "", "Start Date")
	endFlag   = flag.String("EndDate", "", "End Date")

	queries []elastic.Query
)

func init() {
	application.Configure()
}

func main() {
	filter()

	log.Info().Msgf("total mileage: %v", bike.MileageAggregation(application.Context, queries...))
}

func filter() {
	if len(*bikeFlag) > 0 {
		queries = append(queries, bike.Bike(*bikeFlag))
	}

	var start, end *time.Time

	if len(*startFlag) > 0 {
		parsedStart, parseErr := time.Parse(dateFormat, *startFlag)
		logging.Panic(parseErr)

		start = &parsedStart
	}
	if len(*endFlag) > 0 {
		parsedEnd, parseErr := time.Parse(dateFormat, *endFlag)
		logging.Panic(parseErr)

		end = &parsedEnd
	}

	if start != nil || end != nil {
		queries = append(queries, bike.NewDateRangeQuery(start, end))
	}

	if len(queries) == 0 {
		queries = append(queries, bike.All())
	}

	//log.Info().Msgf("queries: %s", queries)
}
