package main

func main() {
  config := loadConfig("config.json")
  conn   := amqpConnect(config.Amqpurl)

  for msg := range merge(amqpFollowQueues(conn, config.Queues)) {
    handleMsg(config.Handlerdir, msg)
  }

}
