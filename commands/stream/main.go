package main

import (
	"github.com/walterjwhite/go-application/libraries/stream"
	"github.com/walterjwhite/go-application/libraries/stream/plugins/database"
	"github.com/walterjwhite/go-application/libraries/stream/plugins/elasticsearch"
	
	"github.com/walterjwhite/go-application/libraries/application"
	elasticsearchl "github.com/walterjwhite/go-application/libraries/elasticsearch"
	databasel "github.com/walterjwhite/go-application/libraries/database"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
	"strconv"
	
	// import all necessary database drivers here ...
)

var (
	esIndexOperationFlag := flag.String("IndexOperation", "Index", "Index Operation")
)

func init() {
	application.Configure()
}

// TODO: 
// 1. define datatype here
// 2. configuration / property (query, username, password)
// 3. how to signal the end of the source
// 4. how to stop processing at the sink once there are no more records
type DataType struct {
	
}

func main() {
	// stream data from source to sink
	sink := &elasticsearch.Sink{NodeConfiguration:
		elasticsearchl.NewDefaultClient(), IndexOperation: bulk.Operation(*esIndexOperationFlag)}
	
	source := &database.Source{Query: databasel.Query{
		QueryString: "",
		Parameters: []string{},
		
		ConnectionConfiguration: databasel.ConnectionConfiguration{
			Username:
			Password:
			Host:
			Port
			Service:
			DriverName: }
		}}
	
	// database source
	// elasticsearch sink
	// datatype
	
	stream.Pipe(source, sink)
}
