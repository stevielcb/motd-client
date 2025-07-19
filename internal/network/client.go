// Package network handles network communication with the MOTD server.
package network

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"
)

// Client handles communication with the MOTD server.
type Client struct {
	host    string
	port    int
	timeout time.Duration
}

// NewClient creates a new network client.
func NewClient(host string, port int, timeout time.Duration) *Client {
	return &Client{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

// Connect establishes a connection to the MOTD server.
func (c *Client) Connect() (net.Conn, error) {
	address := net.JoinHostPort(c.host, fmt.Sprintf("%d", c.port))

	slog.Debug("Connecting to server", "address", address, "timeout", c.timeout)

	conn, err := net.DialTimeout("tcp", address, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	slog.Debug("Successfully connected to server", "address", address)
	return conn, nil
}

// FetchMessage reads the message from the server connection.
func (c *Client) FetchMessage(conn net.Conn) (string, error) {
	var buf bytes.Buffer

	// Set a deadline for reading
	if err := conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return "", fmt.Errorf("failed to set read deadline: %w", err)
	}

	_, err := io.Copy(&buf, conn)
	if err != nil {
		return "", fmt.Errorf("failed to read from connection: %w", err)
	}

	message := buf.String()
	slog.Debug("Message received", "length", len(message))

	return message, nil
}
