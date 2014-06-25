package main

import (
  "io"
  "log"
  "bytes"
  "encoding/json"
)

func handleMsg(msg []byte) {
  brd := bytes.NewReader(msg)
  dec := json.NewDecoder(brd)
  for {
    var m Message
    if err := dec.Decode(&m); err == io.EOF {
      break
    } else if err != nil {
      failOnError(err, "Failed to decode message's json payload")
    }
    log.Printf("Handling with handler '%s'", m.Handler)
  }

}

type Message struct {
  Handler string `json:"handler"`
}
