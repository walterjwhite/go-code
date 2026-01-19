package freeport

import (
	"net"
	"testing"
)

func TestGetRandomUnusedPort(t *testing.T) {
	port, err := GetRandomUnusedPort()
	if err != nil {
		t.Fatalf("GetRandomUnusedPort returned an error: %v", err)
	}

	if port <= 0 {
		t.Errorf("GetRandomUnusedPort returned an invalid port: %d", port)
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		t.Fatalf("Failed to listen on 0.0.0.0:0: %v", err)
	}
	defer close(listener)


}

func TestClose(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create dummy listener: %v", err)
	}

	close(listener)

	close(listener)

	listener2, err := net.Listen("tcp", listener.Addr().String())
	if err == nil {
		err = listener2.Close()
		if err != nil {
			t.Logf("Failed to close listener: %v", err)
		}
	}
}
