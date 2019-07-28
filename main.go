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
	start string
	end   string
)

func main() {
	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		time.Duration(c.TimeoutMs)*time.Millisecond,
	)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	var buf bytes.Buffer
	io.Copy(&buf, conn)

	_, ok := os.LookupEnv("TERM")
	if ok {
		term := os.Getenv("TERM")
		if strings.HasPrefix(term, "screen") {
			start = "\033Ptmux;\033\033]"
			end = "\a\033\\"
		} else {
			start = "\033]"
			end = "\a"
		}
	} else {
		os.Exit(1)
	}

	fmt.Printf("%s%s%s\n", start, &buf, end)
}
