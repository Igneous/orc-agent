package main

import (
  "log"
  "github.com/streadway/amqp"
)

func amqpConnect(url string) (*amqp.Connection) {
  var cfg amqp.Config
  cfg.Properties = amqp.Table {
    "product": "orc-agent",
    "version": "oh-so-alpha",
    }
  conn, err := amqp.DialConfig(url, cfg)
  failOnError(err, "Failed to connect to RabbitMQ")
  // defer conn.Close()
  return conn
}

func amqpSetupChannel(conn *amqp.Connection) (*amqp.Channel) {
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  // defer ch.Close()
  return ch
}

func amqpSetupQueue(ch *amqp.Channel, queue queueconfig) (amqp.Queue) {
  // TODO
  // We need to trap this error, since the queue may already be declared.
  q, err := ch.QueueDeclare(
    queue.Name,       // name
    queue.Durable,    // durable
    queue.AutoDelete, // delete when usused
    queue.Exclusive,  // exclusive
    queue.NoWait,     // noWait
    nil,              // arguments
  )
  failOnError(err, "Failed to declare a queue")

  err = ch.QueueBind(
    q.Name,         // queue name
    queue.Key,      // routing key
    queue.Exchange, // exchange
    queue.NoWait,   // noWait
    nil,            // arguments
  )
  failOnError(err, "Failed to bind a queue")
  return q
}

func amqpRegisterConsumer(ch *amqp.Channel, q queueconfig) <-chan amqp.Delivery {
  msgs, err := ch.Consume(
    q.Name,       // queue name
    "",           // consumer
    true,         // autoAck
    q.Exclusive,  // exclusive
    false,        // noLocal
    q.NoWait,     // noWait
    nil,          // arguments
  )
  failOnError(err, "Failed to register a consumer")
  log.Printf("[amqpRegisterConsumer] Registered consumer for queue: %q", q.Name)
  return msgs
}

func amqpFollowQueue(conn *amqp.Connection, queue queueconfig) <-chan []byte {
  ch   := amqpSetupChannel(conn)
  q    := amqpSetupQueue(ch, queue)
  msgs := amqpRegisterConsumer(ch, queue)
  log.Printf("[amqpFollowQueue] Following %q", q.Name)

  out  := make(chan []byte, 10)
  go func() {
    for d := range msgs {
      out <- d.Body
    }
    close(out)
  }()
  return out
}

func amqpFollowQueues(conn *amqp.Connection, queues []queueconfig) []<-chan[]byte {
  chans := make([]<-chan[]byte, len(queues))
  for i, queue := range queues {
    chans[i] = make(<-chan []byte, 10)
    chans[i] = amqpFollowQueue(conn, queue)
  }
  return chans
}
