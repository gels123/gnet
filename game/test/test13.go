package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 0)
	signal.Notify(ch)

	panic(errors.New("eeeeeerrr")) //panic类似try catch机制

	s := <-ch
	fmt.Println("signal====", s, s.String())
}
