package main

import (
  "net"
  "time"
)

var (
  telnetOut net.Conn
)

func client() {
  // connect to this port
  telnetOut, _ = net.Dial("tcp", dialAddress + ":" + dialPort)
  for {
    buf := make([]byte, 1024)
    l, err := telnetOut.Read(buf)
    if err != nil {
      break
    }
    res := buf[:l]
    // Add incoming traffic 200ms delay for authenticity of dial up connection
    time.Sleep(200 * time.Millisecond)
    baudInBuffer = append(baudInBuffer, bytesToBuffer8N1(res)...)
    readBuffer = append(readBuffer, res...)
  }
  if telnetOut != nil {
    telnetOut.Close()
    telnetOut = nil
  }
  if telnetIn != nil {
    telnetIn.Close()
    telnetIn = nil
  }
}
