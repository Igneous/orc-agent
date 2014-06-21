package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"fmt"
	"net"
	"io/ioutil"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func getFilenames(dir string) ([]string, error) {
  fis, err := ioutil.ReadDir(dir)
  if err != nil {
    return nil, err
  }
  res := make([]string, len(fis))
  for i, fi := range fis {
    res[i] = fi.Name()
  }
  return res, nil
}

func getInterfaceByName(name string) (*net.Interface) {
    inter, err := net.InterfaceByName(name)
    if err != nil {
        panic("Unable to get interface by name")
    }
    return inter
}

func main() {
	conn, err := amqp.Dial("amqp://orc-agent:famc@10.10.10.10:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

  mac := getInterfaceByName("eth0").HardwareAddr.String()

	q, err := ch.QueueDeclare(
		mac,     // name
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
    false,
    nil,
  )
  failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

  handlers, err := getFilenames("/usr/share/orc-agent/handlers")
	failOnError(err, "Could not find any handlers in handlerdir")
	log.Printf("Found handlers: %q", handlers)

	done := make(chan bool)
	go func() {
    counter := 0
		for d := range msgs {
      counter = counter + 1
			log.Printf("[%d] Received a message: %s", counter, d.Body)
			// done <- true
		}
	}()

	log.Printf(" [*] Waiting for messages on %q. To exit press CTRL+C", mac)
	<-done
	log.Printf("Done")

	os.Exit(0)
}
