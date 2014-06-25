package main

import (
  "sync"
)

func merge(cs []<-chan[]byte) <-chan []byte {
    var wg sync.WaitGroup
    out := make(chan[]byte, 10)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan[]byte) {
        for i := range c {
            out <- i
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
