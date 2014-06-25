package main

import (
  "log"
  "github.com/streadway/amqp"
)

func amqpConnect(url string) (*amqp.Connection) {
  conn, err := amqp.Dial("amqp://orc-agent:famc@10.10.10.10:5672/")
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

func amqpSetupQueue(ch *amqp.Channel, queue string) (amqp.Queue) {
  // TODO
  // We need to trap this error, since the queue may already be declared.
  q, err := ch.QueueDeclare(
    queue,   // name
    false,   // durable
    false,   // delete when usused
    false,   // exclusive
    false,   // noWait
    nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")

  err = ch.QueueBind(
    q.Name,       // queue name
    "",           // routing key
    "amq.direct", // exchange
    false,        // noWait
    nil,          // arguments
  )
  failOnError(err, "Failed to bind a queue")
  return q
}

func amqpRegisterConsumer(ch *amqp.Channel, queue amqp.Queue) <-chan amqp.Delivery {
  msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
  failOnError(err, "Failed to register a consumer")
  log.Printf(" [amqpRegisterConsumer] Registered consumer for queue: %q", queue.Name)
  return msgs
}

func amqpFollowQueue(conn *amqp.Connection, queue string) <-chan []byte {
  ch   := amqpSetupChannel(conn)
  q    := amqpSetupQueue(ch, queue)
  msgs := amqpRegisterConsumer(ch, q)
  log.Printf(" [amqpFollowQueue] Following %q", q.Name)

  out  := make(chan []byte, 10)
  go func() {
    for d := range msgs {
      out <- d.Body
    }
    close(out)
  }()
  return out
}

func amqpFollowQueues(conn *amqp.Connection, queues []string) []<-chan[]byte {
  chans := make([]<-chan[]byte, len(queues))
  for i, queue := range queues {
    chans[i] = make(<-chan []byte, 10)
    chans[i] = amqpFollowQueue(conn, queue)
  }
  return chans
}
