package freeport

import (
	"github.com/rs/zerolog/log"
	"net"
)

func GetRandomUnusedPort() (int, error) {
	listener, err := net.Listen("tcp", ":0") // let OS choose a random port
	if err != nil {
		return 0, err
	}

	defer close(listener)
	address := listener.Addr().(*net.TCPAddr)

	return address.Port, nil
}

func close(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		log.Warn().Msgf("failed to close listener: %v", err)
	}
}
