package broker

import (
	"encoding/json"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/timtosi/go-broke/queue"
)

// -----------------------------------------------------------------------------

// Send emits the message `msg` to the first `worker.ZMQWorker` available.
func (b *ZMQBroker) Send(msg string) error {
	identity, _ := b.soc.Recv(0)
	b.soc.Send(identity, zmq.SNDMORE)

	b.soc.Recv(0)
	msgRecv, _ := b.soc.Recv(0)

	msgID, _ := strconv.Atoi(msgRecv)
	b.q.Discard(msgID)

	b.soc.Send("", zmq.SNDMORE)
	_, err := b.soc.Send(msg, 0)
	return err
}

// Close releases resources acquired by `b.soc`.
//
// NOTE: If not called explicitly, the socket will be closed
// on garbage collection.
// NOTE: Address used by `w.soc` could not be available as soon as the function
// returns. See `http://zeromq.org/whitepapers%3aarchitecture#toc6` for details.
func (b *ZMQBroker) Close() { _ = b.soc.Close() }

// -----------------------------------------------------------------------------

// ZMQBroker is a structure representing a message broker. The network
// communication stack lies on a Go implementation of the ZeroMQ library.
type ZMQBroker struct {
	q             *queue.Queue
	soc           *zmq.Socket
	emitAgainChan chan *queue.Message
}

// NewZMQBroker returns a new `broker.ZMQBroker`.
//
// NOTE: `addr` must be of the following form
// - `tcp://<hostname>:<port>` for "regular" TCP networking.
// - `inproc://<name>` for in-process networking.
// - `ipc:///<tmp/filename>` for inter-process communication.
func NewZMQBroker(addr string) (*ZMQBroker, error) {
	soc, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, err
	}
	if err := soc.Bind(addr); err != nil {
		return nil, err
	}
	return &ZMQBroker{emitAgainChan: make(chan *queue.Message, 0), q: queue.NewQueue(), soc: soc}, nil
}

// Run launches the broker and coordinates the message queue `b.q`.
//
// NOTE: This function is an infinite loop.
func (b *ZMQBroker) Run(d time.Duration, workChan chan *queue.Message) {
	go b.q.Poll(b.emitAgainChan, d)

	for {
		select {
		case msg := <-b.emitAgainChan:
			m, _ := json.Marshal(*msg)
			b.q.Push(msg.Copy())
			b.Send(string(m))
		default:
			select {
			case msg := <-workChan:
				m, _ := json.Marshal(*msg)
				b.q.Push(msg)
				b.Send(string(m))
			default:
			}
		}
	}
}
