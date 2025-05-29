// Package main implements a client for fetching and displaying Message of the Day (MOTD)
// from a TCP server. It handles various terminal environments and formatting requirements.
package main

import (
	"bytes"
	"fmt"
	"io"
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
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port)),
		time.Duration(c.TimeoutMs)*time.Millisecond,
	)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	var buf bytes.Buffer
	io.Copy(&buf, conn)

	_, ok := os.LookupEnv("TERM")
	if !ok {
		os.Exit(1)
	}

	term_program, ok := os.LookupEnv("TERM_PROGRAM")
	if ok {
		if strings.HasPrefix(term_program, "iTerm") || strings.HasPrefix(term_program, "vscode") {
			iterm2 = true
		}
	}

	_, ok = os.LookupEnv("SSH_CLIENT")
	if ok {
		sshClient = true
	}

	term := os.Getenv("TERM")

	if !iterm2 && !sshClient && strings.HasPrefix(term, "screen") {
		start = "\033Ptmux;\033\033]"
		end = "\a\033\\"
	} else {
		start = "\033]"
		end = "\a"
	}

	fmt.Printf("%s%s%s\n", start, &buf, end)
}
