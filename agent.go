package main

import (
  "os"
)

func main() {
  //handlers, err := getFilenames("/usr/share/orc-agent/handlers")
  //failOnError(err, "Could not find any handlers in handlerdir")
  //log.Printf("Found handlers: %q", handlers)

  config := loadConfig("config.json")
  conn   := amqpConnect(config.Amqpurl)

  for msg := range merge(amqpFollowQueues(conn, config.Queues)) {
    handleMsg(msg)
  }

  os.Exit(0)
}
