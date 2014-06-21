package main

import (
  "fmt"
  "io/ioutil"
)

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

func main() {
  fmt.Println(getFilenames("."))
}
