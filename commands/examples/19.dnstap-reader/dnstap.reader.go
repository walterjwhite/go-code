package main

import (
	"flag"
	"github.com/dnstap/golang-dnstap"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"time"
	"errors"
	"runtime"
)

var (
	DnsTapFilename = flag.String("DnsTapFilename", "", "DnsTapFilename")
)

func init() {
	application.Configure()
}

func main() {
	if len(*DnsTapFilename) == 0 {
		logging.Panic(errors.New("DnsTapFilename is required"))
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	i, err := dnstap.NewFrameStreamInputFromFilename(*DnsTapFilename)
	logging.Panic(err)

	outputChannel := make(chan []byte, 32)
	o := make(chan bool)

	// make this configurable, user passes in an option(s)
	dnstapProcessor := NewUniqueDomainsProcessor()
	//dnstapProcessor := NewElasticSearchProcessor()
	//dnstapProcessor := NewUniqueResponsesProcessor()

	// make this configurable, specify any number of filters, AND, OR them together ...
	//filter := NewClientFilter("10.30.1.20")
	//filter := NewTimeOfDayFilter("00:00:00 EDT", "23:59:59 EDT")
	filter := &TimeOfDayFilter{Start: TimeOfDay{Hour: 0, Minute: 0}, End: TimeOfDay{Hour: 7, Minute: 0}}

	go i.ReadInto(outputChannel)
	go process(filter, dnstapProcessor, outputChannel, o)
	i.Wait()
	<-o

	log.Info().Msg("@ the end")
}

// the channel isn't closed meaning this for loop never terminates
func process(filter Filter, dnstapProcessor DnstapProcessor, outputChannel chan []byte, o chan bool) {
	log.Info().Msg("Processing file")
	//	defer close(outputChannel)

	dt := &dnstap.Dnstap{}
	for frame := range outputChannel {
		logging.Panic(proto.Unmarshal(frame, dt))

		if filter == nil || filter.Matches(dt) {
			dnstapProcessor.Process(dt)
		} else {
			//log.Println("does NOT match")
		}
	}

	log.Info().Msg("Processed file")

	dnstapProcessor.Flush()
	o <- true
	//close(outputChannel)
}

// goals
// list unique domain names
// list all IP addresses
// find blocked domains
// find all queries by client
