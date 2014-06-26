package main

import (
  "os"
  "os/exec"
  "io"
  "io/ioutil"
  "log"
  "bytes"
  "strings"
  "errors"
  "encoding/json"
)

type Message struct {
  Handler string `json:"handler"`
}

func getHandlers(handlerdir string) []string {
  handlers, err := getFilenames(handlerdir)
  failOnError(err, "Could not find any handlers in handlerdir")
  log.Printf("Found handlers: %q", handlers)
  return handlers
}

func parseMsg(msg []byte) Message {
  brd := bytes.NewReader(msg)
  dec := json.NewDecoder(brd)

  var m Message
  for {
    if err := dec.Decode(&m); err == io.EOF {
      break
    } else if err != nil {
      failOnError(err, "Failed to decode message's json payload")
    }
  }

  return m
}

// Returns a path to the handler if it exists in the handlerdir, else err
func checkHandlerExists(handlerdir string, handler string) (string, error) {
  handlers := getHandlers(handlerdir)
  for _, h := range handlers {
    log.Printf("[checkHandlerExists] checking %q == %q", h, handler)
    if h == handler {
      s := []string{handlerdir, handler}
      sj := strings.Join(s, "")
      log.Printf("[checkHandlerExists] %q == %q is true, returning handlerpath of %q", h, handler, sj)
      log.Printf("I got to the sj return.")
      return sj, nil
    }
  }
  log.Printf("I got to the error return.")
  return "", errors.New("handler doesn't exist")
}

func handleMsg(handlerdir string, msg []byte) {
  parsedmsg := parseMsg(msg)
  h, err := checkHandlerExists(handlerdir, parsedmsg.Handler)
  var output []byte
  if err == nil {
    tmpfile, err := ioutil.TempFile(os.TempDir(), parsedmsg.Handler)
    failOnError(err, "Failed to create a tempfile")
    ioutil.WriteFile(tmpfile.Name(), msg, 0644)
    log.Printf("wrote msg to %q", tmpfile.Name())
    cmd := exec.Command(h, tmpfile.Name())
    output, err = cmd.Output()
    if err != nil {
      log.Fatal(err)
    }
  }
  log.Printf("[handleMsg:MSG] %q", msg)
  log.Printf("[handleMsg:OUT] %q", output)
}
