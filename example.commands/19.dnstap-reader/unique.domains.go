package main

import (
	//"fmt"
	"fmt"
	"github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"net"
	"strings"
)

type UniqueDomains struct {
	DomainStatistics map[string]int
}

type DnstapProcessor interface {
	//Process(output chan []byte)
	Process(dnstapRecord *dnstap.Dnstap)
	Flush()
}

func NewUniqueDomainsProcessor() *UniqueDomains {
	return &UniqueDomains{DomainStatistics: make(map[string]int)}
}

func (p *UniqueDomains) Process(dnstapRecord *dnstap.Dnstap) {
	if dnstapRecord.Message.QueryAddress != nil {
		switch *dnstapRecord.Message.Type {
		case dnstap.Message_CLIENT_QUERY,
			dnstap.Message_RESOLVER_QUERY,
			dnstap.Message_AUTH_QUERY,
			dnstap.Message_FORWARDER_QUERY,
			dnstap.Message_TOOL_QUERY:
			p.processQuery(dnstapRecord)
		default:
			return
		}
	}
}

func (p *UniqueDomains) processQuery(dnstapRecord *dnstap.Dnstap) {
	msg := new(dns.Msg)
	logging.Panic(msg.Unpack(dnstapRecord.Message.QueryMessage))

	if len(msg.Question) > 0 {
		//domain := fmt.Sprint(net.IP(msg.Question[0].Name))
		l := len(msg.Question[0].Name)

		domain := strings.ToLower(msg.Question[0].Name)[:l-1]

		count := p.DomainStatistics[domain]
		count++

		p.DomainStatistics[domain] = count

		fmt.Printf("%v -> %v\n", domain, count)
	}
}

func (p *UniqueDomains) Flush() {
	for domain, count := range p.DomainStatistics {
		fmt.Printf("%v -> %v\n", domain, count)
	}
}
