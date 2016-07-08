package worker

import (
	"encoding/json"
	"strconv"

	zmq "github.com/pebbe/zmq4"
	"github.com/timtosi/go-broke/queue"
)

// -----------------------------------------------------------------------------

// Receive returns a `string` containing a message received on `w.soc` or
// an error.
func (w *ZMQWorker) Receive() (m queue.Message, err error) {
	var msg string

	if _, err = w.soc.Send("", zmq.SNDMORE); err != nil {
		return
	}
	if _, err = w.soc.Send(w.lastMsgID, 0); err != nil {
		return
	}
	if _, err = w.soc.Recv(0); err != nil {
		return
	}
	if msg, err = w.soc.Recv(0); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(msg), &m); err != nil {
		return
	}
	w.lastMsgID = strconv.Itoa(m.ID)
	return
}

// Close releases resources acquired by `w.soc`.
//
// NOTE: If not called explicitly, the socket will be closed
// on garbage collection.
// NOTE: Address used by `w.soc` could not be available as soon as the function
// returns. See `http://zeromq.org/whitepapers%3aarchitecture#toc6` for details.
func (w *ZMQWorker) Close() { _ = w.soc.Close() }

// Identity returns a `string` containing `w.soc`s identity or an `error`.
func (w *ZMQWorker) Identity() (string, error) { return w.soc.GetIdentity() }

// -----------------------------------------------------------------------------

// ZMQWorker is a structure representing a worker process. The network
// communication stack lies on a Go implementation of the ZeroMQ library.
type ZMQWorker struct {
	soc       *zmq.Socket
	lastMsgID string
}

// NewZMQWorker returns a new `ZMQWorker`.
//
// NOTE: `id` should be unique.
// NOTE: `addr` must be of the following form
// - `tcp://<hostname>:<port>` for "regular" TCP networking.
// - `inproc://<name>` for in-process networking.
// - `ipc:///<tmp/filename>` for inter-process communication.
func NewZMQWorker(addr, id string) (*ZMQWorker, error) {
	soc, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		return nil, err
	}
	if err := soc.SetIdentity(id); err != nil {
		return nil, err
	}
	if err := soc.Connect(addr); err != nil {
		return nil, err
	}
	return &ZMQWorker{soc: soc}, nil
}
