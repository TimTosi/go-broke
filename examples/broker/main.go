package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/timtosi/go-broke/broker"
	"github.com/timtosi/go-broke/queue"
)

// -----------------------------------------------------------------------------

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments.\n")
	}
}

// -----------------------------------------------------------------------------

func main() {
	b, err := broker.NewZMQBroker(fmt.Sprintf("tcp://*:%s", os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	workChan := make(chan *queue.Message)
	go b.Run(5*time.Millisecond, workChan)

	for i := 1; ; i++ {
		workChan <- queue.NewMessage(i, "Yo")
		fmt.Println("--------NEW MESSAGE SENT IN WORCHAN--------")
	}
}
