package main

import (
	//"fmt"
	"github.com/dnstap/golang-dnstap"
	//"github.com/miekg/dns"
	//"github.com/walterjwhite/go-application/libraries/logging"
	//"log"
	"net"
	//"strings"
)

type ClientFilter struct {
	Address string
}

type Filter interface {
	Matches(dnstapRecord *dnstap.Dnstap) bool
}

func NewClientFilter(clientAddress string) *ClientFilter {
	return &ClientFilter{Address: clientAddress}
}

func (f *ClientFilter) Matches(dnstapRecord *dnstap.Dnstap) bool {
	if dnstapRecord.Message.QueryAddress != nil {
		queryAddress := net.IP(dnstapRecord.Message.QueryAddress).String()
		return f.Address == queryAddress
	}

	return false
}
