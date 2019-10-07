package main

import (
	"github.com/dnstap/golang-dnstap"
	"github.com/golang/protobuf/proto"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"log"
	//"time"
	"runtime"
)

func main() {
	application.Configure()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// TODO: configure this
	fname := "/tmp/dnstap.example"
	i, err := dnstap.NewFrameStreamInputFromFilename(fname)
	logging.Panic(err)

	outputChannel := make(chan []byte, 32)
	done := make(chan bool)

	// make this configurable, user passes in an option(s)
	//dnstapProcessor := NewUniqueDomainsProcessor()
	dnstapProcessor := NewElasticSearchProcessor()
	//dnstapProcessor := NewUniqueResponsesProcessor()

	// make this configurable, specify any number of filters, AND, OR them together ...
	//filter := NewClientFilter("10.30.1.20")
	//filter := NewTimeOfDayFilter("00:00:00 EDT", "23:59:59 EDT")
	filter := &TimeOfDayFilter{Start: TimeOfDay{Hour: 0, Minute: 0}, End: TimeOfDay{Hour: 7, Minute: 0}}

	go i.ReadInto(outputChannel)
	go process(filter, dnstapProcessor, outputChannel, done)
	i.Wait()

	log.Println("@ the end")
	<-done
	log.Println("now done")
}

func process(filter Filter, dnstapProcessor DnstapProcessor, outputChannel chan []byte, done chan bool) {
	log.Println("Processing file")
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

	log.Println("Processed file")

	dnstapProcessor.Flush()
	close(outputChannel)

	done <- true
}

// goals
// list unique domain names
// list all IP addresses
// find blocked domains
// find all queries by client
