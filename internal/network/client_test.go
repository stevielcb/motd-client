package network

import (
	"net"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	host := "localhost"
	port := 8080
	timeout := 5 * time.Second

	client := NewClient(host, port, timeout)

	if client.host != host {
		t.Errorf("Expected host %s, got %s", host, client.host)
	}
	if client.port != port {
		t.Errorf("Expected port %d, got %d", port, client.port)
	}
	if client.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.timeout)
	}
}

func TestClient_Connect(t *testing.T) {
	// Test with invalid host/port combination
	client := NewClient("invalid-host", 99999, 100*time.Millisecond)

	_, err := client.Connect()
	if err == nil {
		t.Error("Expected error when connecting to invalid host, got nil")
	}
}

func TestClient_FetchMessage(t *testing.T) {
	// Create a test server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	// Get the port the server is listening on
	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	// Start server goroutine
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Send test message
		conn.Write([]byte("Test MOTD Message"))
	}()

	// Create client and connect
	client := NewClient("localhost", port, 1*time.Second)
	conn, err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	// Fetch message
	message, err := client.FetchMessage(conn)
	if err != nil {
		t.Fatalf("Failed to fetch message: %v", err)
	}

	expected := "Test MOTD Message"
	if message != expected {
		t.Errorf("Expected message %q, got %q", expected, message)
	}
}

func TestClient_FetchMessage_Empty(t *testing.T) {
	// Create a test server that sends empty message
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Send empty message
		conn.Write([]byte(""))
	}()

	client := NewClient("localhost", port, 1*time.Second)
	conn, err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	message, err := client.FetchMessage(conn)
	if err != nil {
		t.Fatalf("Failed to fetch message: %v", err)
	}

	if message != "" {
		t.Errorf("Expected empty message, got %q", message)
	}
}

func TestClient_FetchMessage_Timeout(t *testing.T) {
	// Create a test server that doesn't send anything
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Don't send anything, just keep connection open
		time.Sleep(2 * time.Second)
	}()

	client := NewClient("localhost", port, 100*time.Millisecond)
	conn, err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	_, err = client.FetchMessage(conn)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
