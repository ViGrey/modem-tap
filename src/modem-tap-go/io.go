package main

import (
  "bytes"
  "encoding/binary"
  "io/ioutil"
  "log"
  "os"
  "time"
)

var (
  pcmOut []byte
  homeDir = os.Getenv("HOME")
)

func makeWav(data []byte) []byte {
  wavFile := []byte("RIFF")
  subchunk1 := []byte("fmt ")
  subchunk1 = append(subchunk1, []byte{16, 0, 0, 0}...)
  subchunk1 = append(subchunk1, []byte{1, 0}...)
  subchunk1 = append(subchunk1, []byte{1, 0}...)
  subchunk1SampleRate := new(bytes.Buffer)
  binary.Write(subchunk1SampleRate, binary.LittleEndian, uint32(sampleRate))
  subchunk1 = append(subchunk1, subchunk1SampleRate.Bytes()...)
  subchunk1ByteRate := new(bytes.Buffer)
  binary.Write(subchunk1ByteRate, binary.LittleEndian, uint32(sampleRate * 2))
  subchunk1 = append(subchunk1, subchunk1ByteRate.Bytes()...)
  subchunk1 = append(subchunk1, []byte{2, 0}...)
  subchunk1 = append(subchunk1, []byte{16, 0}...)
  subchunk2 := []byte("data")
  subchunk2Len := new(bytes.Buffer)
  binary.Write(subchunk2Len, binary.LittleEndian, uint32(len(data)))
  subchunk2 = append(subchunk2, subchunk2Len.Bytes()...)
  subchunk2 = append(subchunk2, data...)
  subchunksLen := new(bytes.Buffer)
  binary.Write(subchunksLen, binary.LittleEndian, uint32(len(subchunk1) +
                                                         len(subchunk2) + 4))
  wavFile = append(wavFile, subchunksLen.Bytes()...)
  wavFile = append(wavFile, []byte("WAVE")...)
  wavFile = append(wavFile, subchunk1...)
  wavFile = append(wavFile, subchunk2...)
  return wavFile
}

func setupWavFolder() {
  if _, err := os.Stat(homeDir + "/Modem-Tap"); os.IsNotExist(err) {
    err := os.Mkdir(homeDir + "/Modem-Tap", 0755)
    if err != nil {
      log.Println("\x1b[91mmodem-tap: Unable to make dir " + homeDir +
                  "/Modem-Tap\x1b[0m")
      os.Exit(1)
    }
  }
}

func writeWav(data []byte) {
  t := time.Now()
  wavFile := makeWav(data)
  setupWavFolder()
  err := ioutil.WriteFile(homeDir + "/Modem-Tap/" + dialAddress + "-" +
                         dialPort + "-" + t.Format("20060102150405") + ".wav",
                         wavFile, 0644)
  if err != nil {
    panic(err)
  }
}
