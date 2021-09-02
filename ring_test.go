package spool

import (
	"testing"
)

func createFilledRing(len int) *Ring {
	r := NewRing(len)
	for i := 0; i < len; i++ {
		r.Push(i)
	}
	return r
}

func TestNew(t *testing.T) {
	max := 99
	r := NewRing(max)
	if r.max != max {
		t.Errorf("new max was not set max. expected %d, got %d", max, r.max)
	}
	if r.first != nil {
		t.Errorf("new first value was not nil")
	}
	if r.current != nil {
		t.Errorf("new current value was not nil")
	}
	r = NewRing(-1)
	if r != nil {
		t.Errorf("negative new failed to produce nil output")
	}
}

func TestLen(t *testing.T) {
	r := NewRing(1)
	if r.Len() != 0 {
		t.Errorf("new length was not 0")
	}
	r.Push(0)
	if r.Len() != 1 {
		t.Errorf("length was incorrect; expected %d, got %d", 1, r.Len())
	}
	r.Push(0)
	if r.Len() != 2 {
		t.Errorf("length was incorrect; expected %d, got %d", 2, r.Len())
	}
}

func TestAvail(t *testing.T) {
	len := 5
	r := NewRing(len)
	if r.Avail() != len {
		t.Errorf("new avail was incorrect; expected %d, got %d", len, r.Avail())
	}
	for i := 1; i < len; i++ {
		r.Push(i)
		expect := len - i
		if r.Avail() != expect {
			t.Errorf("Avail() was incorrect; expected %d, got %d", expect, r.Avail())
		}
	}
}

func TestPushPop(t *testing.T) {
	limit := 5
	r := createFilledRing(limit)
	if r.Value() != 0 {
		t.Errorf("push incorrectly modified current; expected 0, got %d", r.Value())
	}
	for j := 0; j < limit; j++ {
		val := r.Pop()
		if j != val {
			t.Errorf("Pop() did not return expected value; expected %d, got %d", j, val)
		}
	}

}

func TestPushOverflow(t *testing.T) {
	limit := 5
	r := createFilledRing(limit)
	r.Next()
	testVal := r.Value()
	r.Push(5)
	if r.first.Value != testVal {
		t.Errorf("overlow did not properly pop first element, expected %d, got %d", testVal, r.Value())
	}
	if r.Value() != r.first.Value {
		t.Errorf("overflow did not properly update first, expected %d, got %d", r.Value(), r.first.Value)
	}
}

func TestValue(t *testing.T) {
	testVal := 999
	r := NewRing(1)
	r.Push(testVal)
	val := r.Value()
	if val != testVal {
		t.Errorf("fetched value was incorrect; expected %d, got %d", val, testVal)
	}
}

func TestSetValue(t *testing.T) {
	testVal := 999
	r := NewRing(1)
	r.Push(1)
	r.SetValue(testVal)
	val := r.Value()
	if val != testVal {
		t.Errorf("set value was incorrect; expected %d, got %d", val, testVal)
	}
}

func TestNextPrev(t *testing.T) {
	limit := 5
	r := createFilledRing(limit)
	for j := 0; j < limit; j++ {
		if j != r.Value() {
			t.Errorf("next did not iterate value correctly; expected %d, got %d", j, r.Value())
		}
		r.Next()
	}
	for k := limit - 1; k >= 0; k-- {
		r.Prev()
		if k != r.Value() {
			t.Errorf("prev did not iterate value correctly: expected %d, got %d", k, r.Value())
		}
	}
}

func TestMove(t *testing.T) {
	r := createFilledRing(100)
	r.Move(50)
	if r.Value() != 50 {
		t.Errorf("move did not iterate to proper location; expected %d, got %d", 50, r.Value())
	}
	r.Move(-50)
	if r.Value() != 0 {
		t.Errorf("move did not iterate to proper location; expected %d, got %d", 50, r.Value())
	}
}

func TestDo(t *testing.T) {
	limit := 10
	r := createFilledRing(limit)
	sum := 0
	numIterations := 0
	actualSum := 0
	for i := 0; i < limit; i++ {
		actualSum += i
	}
	r.Do(func(p interface{}) {
		numIterations++
		if p != nil {
			sum += p.(int)
		}
	})
	if numIterations != limit {
		t.Errorf("Do() did not have the proper number of forward iterations, expected %d, got %d", limit, numIterations)
	}
	if actualSum >= 0 && sum != actualSum {
		t.Errorf("forward ring sum = %d, expected %d", sum, actualSum)
	}
}

func BenchmarkPush(b *testing.B) {
	r := NewRing(b.N)
	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	r := createFilledRing(b.N)
	for i := 0; i > b.N; i++ {
		_ = r.Pop()
	}
}
