# go-broke

[![GoDoc](https://godoc.org/github.com/timtosi/go-broke?status.svg)](https://godoc.org/github.com/timtosi/go-broke)
[![Go Report Card](https://goreportcard.com/badge/github.com/timtosi/go-broke)](https://goreportcard.com/report/github.com/timtosi/go-broke)


## Overview
go-broke is a Go package that provides an out-of-the-box message passing tool built on top of [zmq4](https://github.com/pebbe/zmq4),
the latest Golang implementation of the ZeroMQ library by [pebbe](https://github.com/pebbe).

go-broke provides:
* Lightweight Broker object that relies on zmq Router socket.
* Lightweight Worker object that relies on zmq Dealer socket.
* Automatic reemission of messages after a configurable amount of time.

## Getting Started
First of all, you need to install and configure [ZeroMQ](http://zeromq.org/) on your machine.
Check the [install.sh](https://github.com/TimTosi/go-broke/blob/master/install.sh) script for help.

Then go get this repository:

```golang
go get github.com/timtosi/go-broke
```

## Usage
The smallest helloworld of this tool looks like that:

### Broker
```go
func main() {
	b, err := broker.NewZMQBroker("tcp://*:6767")
	if err != nil {
		log.Fatal(err)
	}

	workChan := make(chan *queue.Message)
	go b.Run(50*time.Millisecond, workChan)

	for i := 1; ; i++ {
		workChan <- queue.NewMessage(i, "Yo")
		fmt.Println("--------NEW MESSAGE SENT IN WORCHAN--------")		
	}
}
```

### Worker
```go
func main() {
	w, err := worker.NewZMQWorker("tcp://127.0.0.1:6767", "1")
	if err != nil {
		log.Fatal(err)
	}

	workerID, _ := w.Identity()

	for {
		if msg, err := w.Receive(); err == nil {
			fmt.Printf("Worker %s - Message Received: %v\n", workerID, msg)
		} else {
			fmt.Printf("Worker %s - BUG HERE %v\n", workerID, err)
		}
	}

}
```

## Documentation and Complete Examples
[godoc](https://godoc.org/github.com/timtosi/go-broke)

[example/broker](https://github.com/TimTosi/go-broke/blob/master/examples/broker/main.go)

[example/worker](https://github.com/TimTosi/go-broke/blob/master/examples/worker/main.go)

## License
This library is released under the [MIT License](http://opensource.org/licenses/MIT)

## Not Good Enough ?
Help me to improve by sending your thoughts to timothee.tosi@gmail.com !
