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

type message struct {
  Handler string `json:"handler"`
}

func getHandlers(handlerdir string) []string {
  handlers, err := getFilenames(handlerdir)
  failOnError(err, "Could not find any handlers in handlerdir")
  log.Printf("[getHandlers] Found handlers: %q", handlers)
  return handlers
}

func parseMsg(msg []byte) message {
  brd := bytes.NewReader(msg)
  dec := json.NewDecoder(brd)

  var m message
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
      sj := strings.Join(s, "/")
      log.Printf("[checkHandlerExists] %q == %q is true, returning handlerpath of %q", h, handler, sj)
      return sj, nil
    }
  }
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
    log.Printf("[handleMsg] wrote msg to %q", tmpfile.Name())
    log.Printf("[handleMsg] running to '%s %s'", h, tmpfile.Name())
    cmd := exec.Command(h, tmpfile.Name())
    stdoutpipe, err := cmd.StdoutPipe()
    stderrpipe, err := cmd.StderrPipe()
    go logCmdOutputStream(stdoutpipe)
    go logCmdOutputStream(stderrpipe)
    cmd.Run()
    if err != nil {
      log.Fatal(err)
    }
  }
  log.Printf("[handleMsg:MSG] %s", msg)
  log.Printf("[handleMsg:OUT] %s", output)
}
