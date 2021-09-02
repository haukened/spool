package spool

import (
	"container/ring"
	"sync"
)

// Ring provides a thread safe length limited ring buffer that operates in a FIFO manner
type Ring struct {
	max, used      int
	current, first *ring.Ring
	lock           sync.Mutex
}

func (r *Ring) init(max int) *Ring {
	r.max = max
	return r
}

// New creates a ring with max n elements
// this does not actually create the elements
func NewRing(max int) *Ring {
	if max <= 0 {
		return nil
	}
	r := new(Ring)
	r.init(max)
	return r
}

// Len returns the current length of the ring
func (r *Ring) Len() int {
	return r.current.Len()
}

// Avail returns the current space available in the ring
func (r *Ring) Avail() int {
	return r.max - r.used
}

// Pop removes the oldest element from the ring;
func (r *Ring) Pop() interface{} {
	r.lock.Lock()
	val := r.pop()
	r.lock.Unlock()
	return val
}

func (r *Ring) pop() interface{} {
	r.first = r.first.Prev()
	p := r.first.Unlink(1)
	r.first = r.first.Next()
	return p.Value
}

// Push adds elements to the ring; if the ring is full it replaces the oldest element
func (r *Ring) Push(val interface{}) {
	elem := ring.New(1)
	elem.Value = val
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.first == nil {
		r.current = elem
		r.first = elem
		r.used++
		return
	}
	if r.Avail() <= 0 {
		_ = r.pop()
	}
	r.current = r.current.Prev()
	r.current = r.current.Link(elem)
	r.used++
}

// Value returns the value of the current ring element
func (r *Ring) Value() interface{} {
	return r.current.Value
}

// SetValue changes the value of the current ring element
func (r *Ring) SetValue(val interface{}) {
	r.lock.Lock()
	r.current.Value = val
	r.lock.Unlock()
}

// Next moves the ring one element forward, and returns the destination value. r must not be empty
func (r *Ring) Next() interface{} {
	r.lock.Lock()
	r.current = r.current.Next()
	r.lock.Unlock()
	return r.current.Value
}

// Prev rmoves the ring element backward, and returns the destination value. r must not be empty
func (r *Ring) Prev() interface{} {
	r.lock.Lock()
	r.current = r.current.Prev()
	r.lock.Unlock()
	return r.current.Value
}

// Move moves the currnet position n % r.Len() elements backwards (n < 0) or forward (n > 0)
// in the ring and returns the destination value. r must not be empty
func (r *Ring) Move(n int) interface{} {
	r.lock.Lock()
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r.current = r.current.Prev()
		}
	case n > 0:
		for ; n > 0; n-- {
			r.current = r.current.Next()
		}
	}
	r.lock.Unlock()
	return r.current.Value
}

// Do calls function f on each element in the ring iterating forward from the current position
func (r *Ring) Do(f func(interface{})) {
	r.lock.Lock()
	r.current.Do(f)
	r.lock.Unlock()
}
