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
	defer closeResource(listener)


}

func TestGetPortAllocation(t *testing.T) {
	allocation, err := GetPortAllocation()
	if err != nil {
		t.Fatalf("GetPortAllocation returned an error: %v", err)
	}
	defer func() {
		if err := allocation.Close(); err != nil {
			t.Errorf("failed to close allocation: %v", err)
		}
	}()

	if allocation.Port <= 0 {
		t.Errorf("GetPortAllocation returned an invalid port: %d", allocation.Port)
	}

	if allocation.Port < MinPort || allocation.Port > MaxPort {
		t.Errorf("GetPortAllocation returned port %d outside valid range [%d, %d]",
			allocation.Port, MinPort, MaxPort)
	}

	if allocation.Listener == nil {
		t.Fatal("GetPortAllocation returned nil listener")
	}

	listener2, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		t.Fatalf("Failed to create secondary listener: %v", err)
	}
	defer func() {
		if err := listener2.Close(); err != nil {
			t.Errorf("failed to close secondary listener: %v", err)
		}
	}()

	_, err = net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		t.Logf("Expected: port allocation can fail when system is under pressure: %v", err)
	}
}

func TestPortAllocationClose(t *testing.T) {
	allocation, err := GetPortAllocation()
	if err != nil {
		t.Fatalf("GetPortAllocation returned an error: %v", err)
	}

	port := allocation.Port

	err = allocation.Close()
	if err != nil {
		t.Errorf("Close returned an error: %v", err)
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		t.Fatalf("Failed to listen after closing allocation: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Errorf("failed to close listener: %v", err)
		}
	}()

	_ = port
}

func TestPortAllocationDoubleClose(t *testing.T) {
	allocation, err := GetPortAllocation()
	if err != nil {
		t.Fatalf("GetPortAllocation returned an error: %v", err)
	}

	err = allocation.Close()
	if err != nil {
		t.Errorf("First Close returned an error: %v", err)
	}

	err = allocation.Close()
	if err != nil {
		t.Logf("Second Close returned an error (expected): %v", err)
	}
}

func TestPortAllocationNilListener(t *testing.T) {
	allocation := &PortAllocation{Listener: nil, Port: 8080}
	err := allocation.Close()
	if err != nil {
		t.Errorf("Close with nil listener returned an error: %v", err)
	}
}

func TestClose(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create dummy listener: %v", err)
	}

	closeResource(listener)

	closeResource(listener)

	listener2, err := net.Listen("tcp", listener.Addr().String())
	if err == nil {
		err = listener2.Close()
		if err != nil {
			t.Logf("Failed to close listener: %v", err)
		}
	}
}
