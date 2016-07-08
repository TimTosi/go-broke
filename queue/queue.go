package queue

import (
	"sync"
	"time"
)

// -----------------------------------------------------------------------------

// Shift returns the `queue.Message` contained at the index 0 of `q.msgs` if
// the length of `q.msgs` is greater than 0. It returns `nil` otherwise.
//
// NOTE: This function is thread-safe.
func (q *Queue) Shift() *Message {
	q.Lock()
	defer q.Unlock()

	if len(q.msgs) != 0 {
		msg := q.msgs[0]
		q.msgs = q.msgs[1:]
		return msg
	}
	return nil
}

// Push appends `msg` at the end of `q.msgs`.
//
// NOTE: This function is thread-safe.
func (q *Queue) Push(msg *Message) {
	q.Lock()
	defer q.Unlock()

	if msg == nil {
		return
	}
	q.msgs = append(q.msgs, msg)
}

// Poll purges `q.msgs` from `queue.Message` no longer relevant.
// After `d` time, `q.Purge()` is called and returned value is sent through
// `emitAgainChan` channel if not `nil`.
//
// NOTE: This function is an infinite loop.
func (q *Queue) Poll(emitAgainChan chan<- *Message, d time.Duration) {
	for {
		select {
		case <-time.After(d):
			if msg := q.Purge(d); msg != nil {
				emitAgainChan <- msg
			}
		}
	}
}

// Discard deletes the first `queue.Message` contained in `q.msgs` where
// `queue.Message.ID` matches `ID`.
//
// NOTE: This function is thread-safe.
func (q *Queue) Discard(ID int) bool {
	q.Lock()
	defer q.Unlock()

	for i := 0; i < len(q.msgs); i++ {
		if q.msgs[i].ID == ID {
			copy(q.msgs[i:], q.msgs[i+1:])
			q.msgs[len(q.msgs)-1] = nil
			q.msgs = q.msgs[:len(q.msgs)-1]
			return true
		}
	}
	return false
}

// Purge deletes one `queue.Message` from `q.msgs` where the time elapsed
// since `queue.Message.Timeout` is greater or equal to `d`.
// The `queue.Message` deleted is returned.
//
// NOTE: This function is thread-safe.
func (q *Queue) Purge(d time.Duration) *Message {
	q.Lock()
	defer q.Unlock()

	for i := 0; i < len(q.msgs); i++ {
		if q.msgs[i].TimeoutReached(d) == true {
			msg := q.msgs[i]
			copy(q.msgs[i:], q.msgs[i+1:])
			q.msgs[len(q.msgs)-1] = nil
			q.msgs = q.msgs[:len(q.msgs)-1]
			return msg
		}
	}
	return nil
}

// Queue is a structure representing a first-in first-out container used by
// `queue.ZMQBroker` for re-emitting a given `queue.Message`s in the case
// where a `queue.ZMQWorker` crashes while processing this `queue.Message`.
type Queue struct {
	*sync.RWMutex

	msgs []*Message
}

// NewQueue returns a new `queue.Queue`.
func NewQueue() *Queue { return &Queue{RWMutex: &sync.RWMutex{}, msgs: make([]*Message, 0)} }
