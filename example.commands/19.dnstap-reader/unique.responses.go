package main

import (
	//"fmt"
	"github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
	"github.com/walterjwhite/go-application/libraries/logging"
	"log"
	//"net"
	"strings"
)

type UniqueResponses struct {
	ResponseStatistics map[string]int
}

func NewUniqueResponsesProcessor() *UniqueResponses {
	return &UniqueResponses{ResponseStatistics: make(map[string]int)}
}

func (p *UniqueResponses) Process(dnstapRecord *dnstap.Dnstap) {
	if dnstapRecord.Message.QueryAddress != nil {
		switch *dnstapRecord.Message.Type {
		case dnstap.Message_CLIENT_RESPONSE,
			dnstap.Message_RESOLVER_RESPONSE,
			dnstap.Message_AUTH_RESPONSE,
			dnstap.Message_FORWARDER_RESPONSE,
			dnstap.Message_TOOL_RESPONSE:
			p.processResponse(dnstapRecord)
		default:
			return
		}
	}
}

func (p *UniqueResponses) processResponse(dnstapRecord *dnstap.Dnstap) {
	msg := new(dns.Msg)
	logging.Panic(msg.Unpack(dnstapRecord.Message.ResponseMessage))

	if len(msg.Question) > 0 {
		response := strings.ToLower(msg.String())

		count := p.ResponseStatistics[response]
		count++

		p.ResponseStatistics[response] = count

		log.Printf("%v -> %v\n", response, count)
	}
}

func (p *UniqueResponses) Flush() {
	for response, count := range p.ResponseStatistics {
		log.Printf("%v -> %v\n", response, count)
	}
}
