package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/part"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
	"strings"
)

var (
	filesFlag            = flag.String("Files", "", "Files")
	tagsFlag             = flag.String("Tags", "", "Tags")
	referenceNumbersFlag = flag.String("ReferenceNumbers", "", "Reference Numbers")
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 2 {
		log.Info().Msgf("Adding Part: %v", flag.Args()[0])

		var files, tags, referenceNumbers []string
		if len(*filesFlag) > 0 {
			files = strings.Split(*filesFlag, ",")
		}
		if len(*tagsFlag) > 0 {
			tags = strings.Split(*tagsFlag, ",")
		}
		if len(*referenceNumbersFlag) > 0 {
			referenceNumbers = strings.Split(*referenceNumbersFlag, ",")
		}

		q := part.Add(application.Context, flag.Args()[0], flag.Args()[1], files, tags, referenceNumbers...)
		log.Info().Msgf("Added Part: %v", q)
	} else {
		logging.Panic(fmt.Errorf("No question raised"))
	}
}
