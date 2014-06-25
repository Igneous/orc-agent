package main

import (
	"os"
	"log"
)

func main() {
  //handlers, err := getFilenames("/usr/share/orc-agent/handlers")
	//failOnError(err, "Could not find any handlers in handlerdir")
	//log.Printf("Found handlers: %q", handlers)

  config := loadConfig("config.json")
  conn   := amqpConnect(config.Amqpurl)
  //msgs_tw := amqpFollowQueue(conn, mac_two)

	log.Printf(" [*] Waiting for messages on %q. To exit press CTRL+C", config.Queues)
  for msg := range merge(amqpFollowQueues(conn, config.Queues)) {
    handleMsg(msg)
  }

	os.Exit(0)
}
