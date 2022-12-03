package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 0)
	signal.Notify(ch)

	s := <-ch
	fmt.Println("signal=", s, s.String())
}
