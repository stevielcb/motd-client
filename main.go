package main

import (
  "bytes"
  "fmt"
  "io"
  "net"
  "os"
  "time"
)

func main() {
  conn, err := net.DialTimeout(
    "tcp",
    fmt.Sprintf("%s:%d", c.Host, c.Port),
    time.Duration(c.TimeoutMs) * time.Millisecond,
  )
  if err != nil {
    os.Exit(1)
  }
  defer conn.Close()

  var buf bytes.Buffer
  io.Copy(&buf, conn)
  fmt.Print(&buf)
}
