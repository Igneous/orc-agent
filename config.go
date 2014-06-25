package main

import (
  "os"
  "encoding/json"
)

type Agentconfig struct {
  Queues    []string
  Amqpurl     string
  Stathost    string
  Handlerpath string
}

func loadConfig(configfile string) (Agentconfig) {
  file, _ := os.Open(configfile)
  decoder := json.NewDecoder(file)
  config  := Agentconfig{}
  err     := decoder.Decode(&config)
  failOnError(err, "Failed to decode config")
  return config
}
