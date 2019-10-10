package main

import (
	"fmt"
	"github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
	"github.com/walterjwhite/go-application/libraries/elasticsearch"
	"github.com/walterjwhite/go-application/libraries/logging"
	"log"
	"net"
	"strings"
	"time"
)

type DnsQueryRequest struct {
	Time          *time.Time
	ClientAddress string
	Query         string
}

type DnsQueryResponse struct {
	Time          *time.Time
	ServerAddress string
	Response      string
}

type ElasticSearchProcessor struct {
	ElasticSearchConfiguration *elasticsearch.NodeConfiguration

	batch *elasticsearch.MasterBatch
}

func NewElasticSearchProcessor() *ElasticSearchProcessor {
	es := elasticsearch.NewDefaultClient()
	return &ElasticSearchProcessor{ElasticSearchConfiguration: es, batch: es.NewBatch(10, 2, 5*time.Second, 1)}
}

// processes a single record and internally flushes the batch as needed
func (p *ElasticSearchProcessor) Process(dnstapRecord *dnstap.Dnstap) {
	clientAddress := getClientAddress(dnstapRecord)
	id := getMessageId(dnstapRecord)

	log.Printf("ID: %v\n", id)
	p.batch.Append(id, build(dnstapRecord, clientAddress))
}

func isRequest(dnstapRecord *dnstap.Dnstap) bool {
	switch *dnstapRecord.Message.Type {
	case dnstap.Message_CLIENT_QUERY,
		dnstap.Message_RESOLVER_QUERY,
		dnstap.Message_AUTH_QUERY,
		dnstap.Message_FORWARDER_QUERY,
		dnstap.Message_TOOL_QUERY:
		return true
	default:
		return false
	}

	return false
}

func build(dnstapRecord *dnstap.Dnstap, clientAddress string) {
	if isRequest(dnstapRecord) {
		buildRequest(dnstapRecord, clientAddress)
	}

	return buildResponse(dnstapRecord)
}

func buildRequest(dnstapRecord *dnstap.Dnstap, clientAddress string) *DnsQueryRequest {
	return &DnsQueryRequest{Time: getTime(dnstapRecord.Message.QueryTimeSec), ClientAddress: clientAddress, Query: getQuery(dnstapRecord)}
}

func buildResponse(dnstapRecord *dnstap.Dnstap, clientAddress string) *DnsQueryResponse {
	return &DnsQueryResponse{Time: getTime(dnstapRecord.Message.ResponseTimeSec) /*ServerAddress: clientAddress,*/, Response: getResponse(dnstapRecord)}
}

func getTime(timeSeconds *uint64) *time.Time {
	if timeSeconds != nil {
		return &time.Unix(int64(*timeSeconds), 0)
	}

	return nil
}

func getQuery(dnstapRecord *dnstap.Dnstap) string {
	msg := new(dns.Msg)
	logging.Panic(msg.Unpack(dnstapRecord.Message.QueryMessage))

	if len(msg.Question) > 0 {
		//domain := fmt.Sprint(net.IP(msg.Question[0].Name))
		l := len(msg.Question[0].Name)
		return strings.ToLower(msg.Question[0].Name)[:l-1]
	}

	return ""
}

func getResponse(dnstapRecord *dnstap.Dnstap) string {
	msg := new(dns.Msg)
	logging.Panic(msg.Unpack(dnstapRecord.Message.ResponseMessage))

	if len(msg.Question) > 0 {
		return msg.String()
	}

	return ""
}

func getClientAddress(dnstapRecord *dnstap.Dnstap) string {
	return net.IP(dnstapRecord.Message.QueryAddress).String()
}

func getMessageId(dnstapRecord *dnstap.Dnstap, clientAddress string) string {
	return fmt.Sprintf("%v.%v.%v", dnstapRecord.Message.Type, clientAddress, getTimestamp(dnstapRecord.Message.QueryTimeSec))
}

const timeFormat = "2006.01.02.15.04.05"

func getTimestamp(secs *uint64) string {
	if secs != nil {
		return (time.Unix(int64(*secs), 0).Format(timeFormat))
	}

	return ""
}

// used to flush the last batch
func (p *ElasticSearchProcessor) Flush() {
	p.batch.Flush()
}
