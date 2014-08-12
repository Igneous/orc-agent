package main

import (
  "os"
  "encoding/json"
)

type agentconfig struct {
  Queues   []queueconfig
  Amqpurl    string
  Stathost   string
  Handlerdir string
}

type queueconfig struct {
  Name        string
  Exchange    string
  Key         string
  Durable     bool
  AutoDelete  bool
  Exclusive   bool
  NoWait      bool
}

func loadConfig(configfile string) (agentconfig) {
  file, _ := os.Open(configfile)
  decoder := json.NewDecoder(file)
  config  := agentconfig{}
  err     := decoder.Decode(&config)
  failOnError(err, "Failed to decode config")
  return config
}
