package main

import (
  "os"
  "encoding/json"
)

type agentconfig struct {
  Queues   []string
  Amqpurl    string
  Stathost   string
  Handlerdir string
}

func loadConfig(configfile string) (agentconfig) {
  file, _ := os.Open(configfile)
  decoder := json.NewDecoder(file)
  config  := agentconfig{}
  err     := decoder.Decode(&config)
  failOnError(err, "Failed to decode config")
  return config
}
