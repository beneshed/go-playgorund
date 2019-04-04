package main

import (
	"fmt"
	"github.com/thebenwaters/playground/sideserver/server"
	"os"
	"os/signal"
)

// My goal here was to test running a server in the background while doing other things

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	s := server.NewServer()
	go s.Run()
	fmt.Println("Running")
	<-stop

}
