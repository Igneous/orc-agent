package main

import (
	"github.com/streadway/amqp"
)

func amqpConnect(url string) (*amqp.Connection) {
  conn, err := amqp.Dial("amqp://orc-agent:famc@10.10.10.10:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  // defer conn.Close()
  return conn
}

func amqpGetChannel(conn *amqp.Connection) (*amqp.Channel) {
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  // defer ch.Close()
  return ch
}
