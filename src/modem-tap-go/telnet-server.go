package main

import (
  "fmt"
  "net"
)

var (
  telnetIn net.Conn
)

func server() {
  ln, _ := net.Listen("tcp", "localhost:" + listenPort)
  fmt.Printf("\x1b[94mListening on localhost:" + listenPort + "\x1b[0m\n")
  telnetIn, _ = ln.Accept()
  go client()
  defer ln.Close()
  for {
    buf := make([]byte, 1024)
    l, err := telnetIn.Read(buf)
    if err != nil {
      break
    }
    req := buf[:l]
    writeBuffer = append(writeBuffer, req...)
    baudOutBuffer = append(baudOutBuffer, bytesToBuffer8N1(req)...)
  }
  if telnetIn != nil {
    telnetIn.Close()
    telnetIn = nil
  }
  if telnetOut != nil {
    telnetOut.Close()
    telnetOut = nil
  }
}
