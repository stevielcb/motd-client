package app

import (
	"net"
	"testing"

	"github.com/stevielcb/motd-client/internal/config"
	"github.com/stevielcb/motd-client/internal/terminal"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Host:      "localhost",
		Port:      8080,
		TimeoutMs: 100,
		LogLevel:  "info",
	}

	app := New(cfg)

	if app.cfg != cfg {
		t.Error("Expected config to be set")
	}
	if app.client == nil {
		t.Error("Expected client to be created")
	}
	if app.detector == nil {
		t.Error("Expected detector to be created")
	}
}

func TestApp_Run_NoServer(t *testing.T) {
	cfg := &config.Config{
		Host:      "localhost",
		Port:      99999, // Invalid port
		TimeoutMs: 100,
		LogLevel:  "info",
	}

	app := New(cfg)
	err := app.Run()

	if err == nil {
		t.Error("Expected error when server is not available")
	}
}

func TestApp_Run_WithServer(t *testing.T) {
	// Create a test server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	// Start server goroutine
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte("Test MOTD Message"))
	}()

	cfg := &config.Config{
		Host:      "localhost",
		Port:      port,
		TimeoutMs: 1000,
		LogLevel:  "info",
	}

	app := New(cfg)
	err = app.Run()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestApp_displayMessage(t *testing.T) {
	// Create a mock environment
	env := &terminal.Environment{
		StartSeq: "\033]",
		EndSeq:   "\a",
	}

	cfg := &config.Config{
		Host:      "localhost",
		Port:      8080,
		TimeoutMs: 100,
		LogLevel:  "info",
	}

	app := New(cfg)
	app.formatter = terminal.NewFormatter(env)

	// Test with non-empty message
	app.displayMessage("Test Message")

	// Test with empty message
	app.displayMessage("")
}

// Mock implementations for testing
type mockDetector struct {
	env *terminal.Environment
	err error
}

func (m *mockDetector) Detect() (*terminal.Environment, error) {
	return m.env, m.err
}

type mockClient struct {
	conn net.Conn
	err  error
}

func (m *mockClient) Connect() (net.Conn, error) {
	return m.conn, m.err
}

func (m *mockClient) FetchMessage(conn net.Conn) (string, error) {
	return "Mock Message", nil
}

func TestApp_Run_WithMocks(t *testing.T) {
	// Create a test server to get a real connection
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	// Start server goroutine
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte("Test Message"))
	}()

	// Connect to the test server to get a real connection
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	cfg := &config.Config{
		Host:      "localhost",
		Port:      8080,
		TimeoutMs: 100,
		LogLevel:  "info",
	}

	app := New(cfg)

	// Test with successful detection
	mockEnv := &terminal.Environment{
		StartSeq: "\033]",
		EndSeq:   "\a",
	}
	app.detector = &mockDetector{env: mockEnv, err: nil}

	// Mock the client with a real connection
	app.client = &mockClient{conn: conn, err: nil}

	// This should work with the mock detector
	err = app.Run()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestApp_Run_DetectionError(t *testing.T) {
	cfg := &config.Config{
		Host:      "localhost",
		Port:      8080,
		TimeoutMs: 100,
		LogLevel:  "info",
	}

	app := New(cfg)

	// Test with detection error
	app.detector = &mockDetector{env: nil, err: terminal.ErrTerminalNotSet}

	err := app.Run()
	if err == nil {
		t.Error("Expected error when terminal detection fails")
	}
}
