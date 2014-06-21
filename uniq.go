package main

import (
	"net"
)

func getInterfaceByName(name string) (*net.Interface) {
  inter, err := net.InterfaceByName(name)
  if err != nil {
      panic("Unable to get interface by name")
  }
  return inter
}
