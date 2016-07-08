package queue

import (
	"testing"
	"time"

	"github.com/facebookgo/ensure"
)

// -----------------------------------------------------------------------------

func TestMessage_TimeoutReached_true(t *testing.T) {
	m := NewMessage(1, "Ok")
	time.Sleep(1 * time.Second)
	ensure.True(t, m.TimeoutReached(1*time.Second))
}

func TestMessage_TimeoutReached_false(t *testing.T) {
	m := NewMessage(1, "Ok")
	ensure.False(t, m.TimeoutReached(1*time.Second))
}

func TestMessage_Copy(t *testing.T) {
	m1 := NewMessage(1, "Ok")
	m2 := m1.Copy()
	ensure.DeepEqual(t, m1.ID, m2.ID)
	ensure.DeepEqual(t, m1.Msg, m2.Msg)
	ensure.NotDeepEqual(t, m1.Timeout, m2.Timeout)
}
