package freeport

import (
	"github.com/rs/zerolog/log"
	"net"
)

func GetRandomUnusedPort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	defer closeResource(listener)
	address := listener.Addr().(*net.TCPAddr)

	return address.Port, nil
}

func closeResource(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		log.Warn().Msgf("failed to close listener: %v", err)
	}
}
