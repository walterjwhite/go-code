package main

import (
	"time"
	"fmt"
	"github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/elasticsearch"
	"net"
	"strings"
)

type ElasticSearchProcessor struct {
	ElasticSearchConfiguration *elasticsearch.NodeConfiguration
	
	batch *elasticsearch.MasterBatch
}

func NewElasticSearchProcessor() *ElasticSearchProcessor {
	es := elasticsearch.NewDefaultClient()
	return &ElasticSearchProcessor{ElasticSearchConfiguration: es, batch: es.NewBatch(1000, 2, 30*time.Second, 2)}
}

// processes a single record and internally flushes the batch as needed
func (p *ElasticSearchProcessor) Process(dnstapRecord *dnstap.Dnstap) {
	id := getMessageId(dnstapRecord)
	p.batch.Append(id, dnstapRecord)
}

func getMessageId(dnstapRecord *dnstap.Dnstap) string {
	return fmt.Sprintf("%v.%v.%v", dnstapRecord.Message.Type,
		net.IP(dnstapRecord.Message.QueryAddress).String(),
		getTimestamp(dnstapRecord.Message.QueryTimeSec)
}

const timeFormat = "2006.01.02.15.04.05"
func getTimestamp(secs *uint64) string {
	if secs != nil {
		return(time.Unix(int64(*secs), 0).Format(timeFormat))
	}
	
	return ""
}

// used to flush the last batch
func (p *ElasticSearchProcessor) Flush() {
	p.batch.Flush()
}
