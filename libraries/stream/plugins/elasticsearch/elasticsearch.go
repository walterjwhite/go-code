package elasticsearch

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/document"
	elasticsearchl "github.com/walterjwhite/go-application/libraries/elasticsearch"
	"github.com/walterjwhite/go-application/libraries/elasticsearch/bulk"
	"github.com/walterjwhite/go-application/libraries/logging"
)

type Sink struct {
	NodeConfiguration *elasticsearchl.NodeConfiguration
	batch             *bulk.MasterBatch

	Operation bulk.Operation
}

func (s *Sink) Read(channel chan interface{}) {
	if s.batch == nil {
		s.batch = bulk.NewDefaultBatch(s.NodeConfiguration)
	}

	// read from the channel
	for {
		s.process(<-channel)
	}
}

func (s *Sink) process(d interface{}) {
	document, ok := d.(document.Document)
	if !ok {
		logging.Panic(fmt.Errorf("Error converting to Document: %v", d))
	}

	if s.Operation == bulk.Index {
		s.batch.Index(document)
	} else if s.Operation == bulk.Update {
		s.batch.Update(document)
	} else {
		s.batch.Delete(document)
	}
}
