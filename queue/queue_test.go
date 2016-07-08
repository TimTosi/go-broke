package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
)

// -----------------------------------------------------------------------------

var (
	funcMap = map[int]func(*Queue){
		0: push,
		1: shift,
		2: discard,
		3: purge,
	}
)

// -----------------------------------------------------------------------------

// NOTE: the following functions are implemented to test concurrency.
func push(q *Queue)    { q.Push(NewMessage(1, "OK")) }
func shift(q *Queue)   { q.Shift() }
func discard(q *Queue) { q.Discard(1) }
func purge(q *Queue)   { q.Purge(10 * time.Millisecond) }

// -----------------------------------------------------------------------------

func TestQueue_Push(t *testing.T) {
	q := NewQueue()
	q.Push(NewMessage(1, "Ok"))
	ensure.Subset(t, q.msgs[0], NewMessage(1, "Ok"))
}

func TestQueue_Push_nil(t *testing.T) {
	q := NewQueue()
	q.Push(nil)
	ensure.DeepEqual(t, len(q.msgs), 0)
}

func TestQueue_Shift(t *testing.T) {
	q := NewQueue()
	q.Push(NewMessage(1, "Ok"))
	ensure.Subset(t, q.Shift(), NewMessage(1, "Ok"))
}

func TestQueue_Shift_emptyQueue(t *testing.T) {
	q := NewQueue()
	ensure.DeepEqual(t, q.Shift(), (*Message)(nil))
}

func TestQueue_Discard(t *testing.T) {
	q := NewQueue()
	msg1 := NewMessage(60, "Discard me!")
	msg2 := NewMessage(17, "I'm good here.")

	q.Push(msg1)
	q.Push(msg2)
	ensure.SameElements(t, q.msgs, []*Message{msg1, msg2})
	ensure.True(t, q.Discard(60))
	ensure.SameElements(t, q.msgs, []*Message{msg2})
}

func TestQueue_Discard_notFound(t *testing.T) {
	q := NewQueue()
	msg1 := NewMessage(60, "Discard me!")
	msg2 := NewMessage(17, "I'm good here.")

	q.Push(msg1)
	q.Push(msg2)
	ensure.SameElements(t, q.msgs, []*Message{msg1, msg2})
	ensure.False(t, q.Discard(777))
	ensure.SameElements(t, q.msgs, []*Message{msg1, msg2})
}

func TestQueue_Discard_empty(t *testing.T) {
	q := NewQueue()
	ensure.False(t, q.Discard(777))
}

func TestQueue_Purge_oneMatch(t *testing.T) {
	q := NewQueue()
	msg1 := NewMessage(60, "Discard me!")
	time.Sleep(1 * time.Second)
	msg2 := NewMessage(17, "I'm good here.")

	q.Push(msg1)
	q.Push(msg2)
	ensure.SameElements(t, q.msgs, []*Message{msg1, msg2})
	ensure.DeepEqual(t, q.Purge(1*time.Second), msg1)
	ensure.SameElements(t, q.msgs, []*Message{msg2})
	ensure.DeepEqual(t, q.Purge(1*time.Second), (*Message)(nil))
}

func TestQueue_Purge_noMatch(t *testing.T) {
	q := NewQueue()
	msg1 := NewMessage(60, "Discard me!")

	q.Push(msg1)
	ensure.DeepEqual(t, q.Purge(1*time.Second), (*Message)(nil))
	ensure.SameElements(t, q.msgs, []*Message{msg1})
}

func TestQueue_Purge_multipleMatches(t *testing.T) {
	q := NewQueue()
	msg1 := NewMessage(60, "Discard me!")
	msg2 := NewMessage(61, "Discard me Too!")
	time.Sleep(1 * time.Second)
	msg3 := NewMessage(17, "I'm good here.")

	q.Push(msg1)
	q.Push(msg2)
	q.Push(msg3)
	ensure.SameElements(t, q.msgs, []*Message{msg1, msg2, msg3})
	ensure.DeepEqual(t, q.Purge(1*time.Second), msg1)
	ensure.SameElements(t, q.msgs, []*Message{msg2, msg3})
	ensure.DeepEqual(t, q.Purge(1*time.Second), msg2)
	ensure.SameElements(t, q.msgs, []*Message{msg3})
	ensure.DeepEqual(t, q.Purge(1*time.Second), (*Message)(nil))
}

func TestQueue_Purge_empty(t *testing.T) {
	q := NewQueue()
	ensure.DeepEqual(t, q.Purge(1*time.Second), (*Message)(nil))
}

// -----------------------------------------------------------------------------

// NOTE: run with `-race`.
func TestQueue_races(t *testing.T) {
	wg := &sync.WaitGroup{}
	q := NewQueue()
	for i := 0; i < 8191; i++ {
		wg.Add(1)
		go func(i int) {
			funcMap[rand.Intn(4)](q)
			wg.Done()
			if i%1000 == 0 {
				fmt.Printf("Loop %d\n", i)
			}
		}(i)
	}
	wg.Wait()
}
