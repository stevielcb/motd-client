// Package main implements a client for fetching and displaying Message of the Day (MOTD)
// from a TCP server. It handles various terminal environments and formatting requirements.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"
)

// Global variables for terminal-specific formatting
var (
	start     string // ANSI escape sequence start
	end       string // ANSI escape sequence end
	iterm2    bool   // Flag for iTerm2/VSCode terminal
	sshClient bool   // Flag for SSH client
)

// main is the entry point of the application. It:
// 1. Connects to the MOTD server
// 2. Fetches the message
// 3. Determines the terminal environment
// 4. Formats and displays the message with appropriate terminal escape sequences
func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

// run contains the main application logic with proper error handling.
func run() error {
	// Connect to the MOTD server
	conn, err := connectToServer()
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Fetch the message
	message, err := fetchMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to fetch message: %w", err)
	}

	// Determine terminal environment
	determineTerminalEnvironment()

	// Format and display the message
	displayMessage(message)

	return nil
}

// connectToServer establishes a connection to the MOTD server.
func connectToServer() (net.Conn, error) {
	address := net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))

	slog.Debug("Connecting to server", "address", address, "timeout", c.Timeout())

	conn, err := net.DialTimeout("tcp", address, c.Timeout())
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	slog.Debug("Successfully connected to server", "address", address)
	return conn, nil
}

// fetchMessage reads the message from the server connection.
func fetchMessage(conn net.Conn) (string, error) {
	var buf bytes.Buffer

	// Set a deadline for reading
	if err := conn.SetReadDeadline(time.Now().Add(c.Timeout())); err != nil {
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

// determineTerminalEnvironment detects the terminal type and sets appropriate formatting.
func determineTerminalEnvironment() {
	// Check if TERM environment variable exists
	_, ok := os.LookupEnv("TERM")
	if !ok {
		slog.Error("TERM environment variable not set")
		os.Exit(1)
	}

	// Check for iTerm2/VSCode terminal
	termProgram, ok := os.LookupEnv("TERM_PROGRAM")
	if ok {
		if strings.HasPrefix(termProgram, "iTerm") || strings.HasPrefix(termProgram, "vscode") {
			iterm2 = true
			slog.Debug("Detected iTerm2/VSCode terminal")
		}
	}

	// Check for SSH client
	_, ok = os.LookupEnv("SSH_CLIENT")
	if ok {
		sshClient = true
		slog.Debug("Detected SSH client")
	}

	// Set terminal-specific formatting
	term := os.Getenv("TERM")
	if !iterm2 && !sshClient && strings.HasPrefix(term, "screen") {
		start = "\033Ptmux;\033\033]"
		end = "\a\033\\"
		slog.Debug("Using tmux-specific formatting")
	} else {
		start = "\033]"
		end = "\a"
		slog.Debug("Using standard terminal formatting")
	}
}

// displayMessage formats and displays the MOTD message.
func displayMessage(message string) {
	if message == "" {
		slog.Warn("Received empty message from server")
		return
	}

	formattedMessage := fmt.Sprintf("%s%s%s", start, message, end)
	fmt.Print(formattedMessage)

	slog.Debug("Message displayed successfully", "message_length", len(message))
}
