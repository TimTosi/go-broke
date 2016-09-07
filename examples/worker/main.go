package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/timtosi/go-broke/worker"
)

// -----------------------------------------------------------------------------

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments.\n")
	}
}

// -----------------------------------------------------------------------------

func main() {
	wq1, err := worker.NewZMQWorker(fmt.Sprintf("tcp://127.0.0.1:%s", os.Args[1]), "1")
	if err != nil {
		log.Fatal(err)
	}
	wq2, err := worker.NewZMQWorker(fmt.Sprintf("tcp://127.0.0.1:%s", os.Args[1]), "2")
	if err != nil {
		log.Fatal(err)
	}
	wq3, err := worker.NewZMQWorker(fmt.Sprintf("tcp://127.0.0.1:%s", os.Args[1]), "3")
	if err != nil {
		log.Fatal(err)
	}

	do := func(w *worker.ZMQWorker) {
		for {
			workerID, _ := w.Identity()
			if msg, err := w.Receive(); err == nil {
				fmt.Printf("Worker %s - Message Received: %v\n", workerID, msg)
			} else {
				fmt.Printf("Worker %s - BUG HERE %v\n", workerID, err)
			}
			time.Sleep(2 * time.Second)
		}
	}

	go do(wq1)
	go do(wq2)
	go do(wq3)

	for {
		time.Sleep(10 * time.Millisecond)
	}
}
