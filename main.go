package main

import (
  "net"
  "os"
  "fmt"
  "bytes"
  "io"
  "time"
)

func main() {
  conn, err := net.DialTimeout("tcp", "localhost:4200", 100 * time.Millisecond)
  if err != nil {
    os.Exit(1)
  }
  defer conn.Close()
  var buf bytes.Buffer
  io.Copy(&buf, conn)
  fmt.Print(&buf)
}
