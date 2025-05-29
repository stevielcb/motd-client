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

var (
	start     string
	end       string
	iterm2    bool
	sshClient bool
)

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
