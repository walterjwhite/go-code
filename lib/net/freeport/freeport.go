package freeport

import (
	"errors"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
)

const (
	MinPort = 1024
	MaxPort = 65535
)

var ErrPortOutOfRange = errors.New("allocated port is out of valid range")

type PortAllocation struct {
	Listener net.Listener
	Port     int
}

func (p *PortAllocation) Close() error {
	if p.Listener != nil {
		return p.Listener.Close()
	}
	return nil
}

func GetRandomUnusedPort() (int, error) {
	allocation, err := GetPortAllocation()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := allocation.Close(); err != nil {
			log.Warn().Msgf("failed to close port allocation: %v", err)
		}
	}()
	return allocation.Port, nil
}

func GetPortAllocation() (*PortAllocation, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to allocate port: %w", err)
	}

	address, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		closeResource(listener)
		return nil, fmt.Errorf("failed to get TCP address from listener")
	}

	port := address.Port

	if port < MinPort || port > MaxPort {
		closeResource(listener)
		return nil, fmt.Errorf("%w: port %d is not in range [%d, %d]", ErrPortOutOfRange, port, MinPort, MaxPort)
	}

	return &PortAllocation{
		Listener: listener,
		Port:     port,
	}, nil
}

func closeResource(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		log.Warn().Msgf("failed to close listener: %v", err)
	}
}
