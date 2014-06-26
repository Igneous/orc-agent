package main

import "flag"

func main() {
  configfile := flag.String("config", "/etc/orc-agent/config.json",
                            "path to orc-agent's json config file")
  flag.Parse()

  config := loadConfig(*configfile)
  conn   := amqpConnect(config.Amqpurl)

  for msg := range merge(amqpFollowQueues(conn, config.Queues)) {
    handleMsg(config.Handlerdir, msg)
  }
}
