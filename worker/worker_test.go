package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	zmq "github.com/pebbe/zmq4"
	"github.com/timtosi/go-broke/queue"
)

// -----------------------------------------------------------------------------

const dummyServerAddr = "tcp://127.0.0.1:3579"

var dummy *dummyServer

// -----------------------------------------------------------------------------

func init() {
	var err error

	dummy, err = newDummyServer()
	if err != nil {
		log.Fatal(err)
	}
}

// -----------------------------------------------------------------------------

// connect is a dummy function for testing purposes.
func (ds *dummyServer) connect(addr string) (err error) {
	ds.soc.Close()
	ds.soc, err = zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if err = ds.soc.Bind(addr); err == nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	return
}

// diconnect is a dummy function for testing purposes.
func (ds *dummyServer) diconnect() { ds.soc.Close() }

// send is a dummy function for testing purposes.
func (ds *dummyServer) send(msg string) error {
	identity, _ := ds.soc.Recv(0)
	ds.soc.Send(identity, zmq.SNDMORE)

	ds.soc.Recv(0)
	ds.soc.Recv(0)

	ds.soc.Send("", zmq.SNDMORE)
	_, err := ds.soc.Send(msg, 0)
	return err
}

// newDummyServer is a dummy function for testing purposes.
func newDummyServer() (*dummyServer, error) {
	soc, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, err
	}
	return &dummyServer{soc: soc}, nil
}

// dummyServer is a dummy structure for testing purposes.
type dummyServer struct {
	soc *zmq.Socket
}

// -----------------------------------------------------------------------------

func TestWorker_NewZMQWorker_singleWorker(t *testing.T) {
	ensure.Nil(t, dummy.connect(dummyServerAddr))
	defer dummy.diconnect()

	w, err := NewZMQWorker(dummyServerAddr, "0")
	ensure.Nil(t, err)
	w.Close()
}

func TestWorker_NewZMQWorker_multipleWorkers(t *testing.T) {
	var workerSlice []*ZMQWorker

	ensure.Nil(t, dummy.connect(dummyServerAddr))
	defer func() {
		dummy.diconnect()
		for i := 0; i < len(workerSlice); i++ {
			workerSlice[i].Close()
		}
	}()

	for i := 0; i < 200; i++ {
		w, err := NewZMQWorker(dummyServerAddr, fmt.Sprintf("%d", i))
		ensure.Nil(t, err)
		workerSlice = append(workerSlice, w)
	}
}

func TestWorker_NewZMQWorker_noServer(t *testing.T) {
	dummy.diconnect()

	_, err := NewZMQWorker(dummyServerAddr, "0")
	ensure.Nil(t, err)
}

func TestWorker_Receive_singleWorker(t *testing.T) {
	ensure.Nil(t, dummy.connect(dummyServerAddr))
	defer dummy.diconnect()

	w, err := NewZMQWorker(dummyServerAddr, "0")
	ensure.Nil(t, err)

	m, _ := json.Marshal(*queue.NewMessage(1, "Yo !"))
	go dummy.send(string(m))
	msg, err := w.Receive()
	ensure.Nil(t, err)

	ensure.DeepEqual(t, msg.ID, 1)
	ensure.DeepEqual(t, msg.Msg, "Yo !")

	w.Close()
}

func TestWorker_Receive_multipleWorkers(t *testing.T) {
	ensure.Nil(t, dummy.connect(dummyServerAddr))
	defer dummy.diconnect()

	w1, err := NewZMQWorker(dummyServerAddr, "1")
	ensure.Nil(t, err)
	w2, err := NewZMQWorker(dummyServerAddr, "2")
	ensure.Nil(t, err)

	m1, _ := json.Marshal(*queue.NewMessage(1, "Yo !"))
	m2, _ := json.Marshal(*queue.NewMessage(2, "Hi !"))
	go func() {
		dummy.send(string(m1))
		dummy.send(string(m2))
	}()
	msg1, err := w1.Receive()
	ensure.Nil(t, err)
	msg2, err := w2.Receive()
	ensure.Nil(t, err)

	ensure.DeepEqual(t, msg1.ID, 1)
	ensure.DeepEqual(t, msg1.Msg, "Yo !")
	ensure.DeepEqual(t, msg2.ID, 2)
	ensure.DeepEqual(t, msg2.Msg, "Hi !")

	w1.Close()
	w2.Close()
}
