package main

import (
  "github.com/gordonklaus/portaudio"
  "bytes"
  "encoding/binary"
  "fmt"
  "math"
  "os"
  "strconv"
)

const sampleRate = 44100
const pcmScale = 32768
const baud = 147
const scale = 0.5

var (
  baudOutBuffer, baudInBuffer []bool
  writeBuffer, readBuffer []byte
  baudOutTick, baudInTick, writeTick, readTick int
  dialAddress, dialPort string
  listenPort = "2600"
  quietFlag, wavFlag bool
)

type sine struct {
  *portaudio.Stream
  phaseOut, phaseIn float64
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func bytesToBuffer8N1(word []byte) (buffer []bool) {
  for _, n := range word {
    a := strconv.FormatInt(int64(n), 2)
    tmpA := a
    for i := 0; i < 8 - len(a); i++ {
      tmpA = "0" + tmpA
    }
    a = tmpA
    buffer = append(buffer, false)
    for _, b := range Reverse(a) {
      if b == '1' {
        buffer = append(buffer, true)
      } else {
        buffer = append(buffer, false)
      }
    }
    buffer = append(buffer, true)
  }
  return
}

func newSine(sampleRate float64) *sine {
  s := &sine{nil, 0, 0}
  var err error
  s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
  chk(err)
  return s
}

func (g *sine) processAudio(out []float32) {
  for i := range out {
    var freqOut float64 = 1270
    var freqIn float64 = 2225
    if telnetIn == nil {
      baudOutTick = 0
      writeTick = 0
      baudOutBuffer = []bool{}
      writeBuffer = []byte{}
      freqOut = 0
    }
    if telnetOut == nil {
      baudInTick = 0
      readTick = 0
      baudInBuffer = []bool{}
      readBuffer = []byte{}
      freqIn = 0
    }
    bitOut := true
    bitIn := true
    if len(baudOutBuffer) > 0 {
      if baudOutTick != baud {
        baudOutTick++
        if baudOutTick == baud {
          if len(baudOutBuffer) > 0 {
            baudOutBuffer = baudOutBuffer[1:]
            writeTick++
            if writeTick == 10 {
              writeTick = 0
              if len(writeBuffer) > 0 {
                telnetOut.Write([]byte{writeBuffer[0]})
                writeBuffer = writeBuffer[1:]
              }
            }
          }
        }
      } else {
        baudOutTick = 0
      }
      if len(baudOutBuffer) > 0 {
        if baudOutBuffer[0] == false {
          bitOut = false
        }
      }
    }
    if len(baudInBuffer) > 0 {
      if baudInTick != baud {
        baudInTick++
        if baudInTick == baud {
          if len(baudInBuffer) > 0 {
            baudInBuffer = baudInBuffer[1:]
            readTick++
            if readTick == 10 {
              readTick = 0
              if len(readBuffer) > 0 {
                telnetIn.Write([]byte{readBuffer[0]})
                readBuffer = readBuffer[1:]
              }
            }
          }
        }
      } else {
        baudInTick = 0
      }
      if len(baudInBuffer) > 0 {
        if baudInBuffer[0] == false {
          bitIn = false
        }
      }
    }
    if !bitOut {
      freqOut = 1070
    }
    if !bitIn {
      freqIn = 2025
    }
    stepOut := freqOut / sampleRate
    stepIn := freqIn / sampleRate
    tmpOut := float32(math.Sin(2 * math.Pi * g.phaseOut) * scale +
                      math.Sin(2 * math.Pi * g.phaseIn) * scale)
    if quietFlag {
      out[i] = 0
    } else{
      out[i] = tmpOut
    }
    if freqOut != 0 && freqIn != 0 && wavFlag {
      pcmBuf16 := new(bytes.Buffer)
      binary.Write(pcmBuf16, binary.LittleEndian, int16(tmpOut * pcmScale))
      pcmOut = append(pcmOut, pcmBuf16.Bytes()...)
    }
    if freqOut == 0 && freqIn == 0 && wavFlag {
      if len(pcmOut) > 0 {
        pcmOutCopy := pcmOut
        pcmOut = []byte{}
        writeWav(pcmOutCopy)
      }
    }
    _, g.phaseOut = math.Modf(g.phaseOut + stepOut)
    _, g.phaseIn = math.Modf(g.phaseIn + stepIn)
  }
}

func chk(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
  if len(os.Args) > 1 {
    for i, arg := range os.Args[1:] {
      if arg == "-h" || arg == "--help" {
        fmt.Println("modem-tap [ OPTIONS ]... [ LISTENING PORT ]\n\n" +
                    "Options:\n\n" +
                    "-h, --help      Print help (this message) and exit\n" +
                    "-q, --quiet     Does not play connection sounds from " +
                    "the speaker\n" +
                    "-w, --wav       Saves connection sounds to a WAV file\n\n" +
                    "LISTENING PORT is 2600 by default")
        os.Exit(0)
      } else if arg == "-w" || arg == "--wav" {
        wavFlag = true
      } else if arg == "-q" || arg == "--quiet"{
        quietFlag = true
      } else if i + 1 == len(os.Args) - 1 {
        listenPort = arg
      }
    }
  }
  portaudio.Initialize()
  defer portaudio.Terminate()
  s := newSine(sampleRate)
  defer s.Close()
  chk(s.Start())
  for {
    var tmpDialAddress, tmpDialPort string
    fmt.Printf("\x1b[92mEnter Server Address:\x1b[0m ")
    fmt.Scanln(&tmpDialAddress)
    dialAddress = tmpDialAddress
    fmt.Printf("\x1b[92mEnter Port [23]:\x1b[0m ")
    fmt.Scanln(&tmpDialPort)
    dialPort = tmpDialPort
    if dialPort == "" {
      dialPort = "23"
    }
    server()
  }
  chk(s.Stop())
}

