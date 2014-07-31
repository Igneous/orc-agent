package main

import (
  "flag"
  "log"
  "os"
)

func main() {
  // Parse CLI options
  configfile := flag.String("config", "/etc/orc-agent/config.json",
                            "path to orc-agent's json config file")
  logfile    := flag.String("logfile", "/var/log/orc-agent.log",
                            "where orc-agent should log to")
  flag.Parse()

  // Set up our logger
  f, _ := os.Create(*logfile)
  log.SetOutput(f)
  defer f.Close() // clean up after ourselves when we exit.

  // Load our config
  config := loadConfig(*configfile)

  // Connect to rabbitmq broker
  conn   := amqpConnect(config.Amqpurl)

  // Subscribe to queues, and call handleMsg for all incoming messages
  for msg := range merge(amqpFollowQueues(conn, config.Queues)) {
    handleMsg(config.Handlerdir, msg)
  }
}
