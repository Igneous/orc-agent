package main

import (
  "bufio"
  "log"
  "io"
)

func logCmdOutputStream(stream io.ReadCloser) {
  scanner := bufio.NewScanner(stream)
  for scanner.Scan() {
    log.Println(scanner.Text())
  }
  defer stream.Close()
}
