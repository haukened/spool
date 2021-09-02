[![Go Report Card](https://goreportcard.com/badge/haukened/fifo)](https://goreportcard.com/report/haukened/fifo) [![Go Reference](https://pkg.go.dev/badge/github.com/haukened/fifo.svg)](https://pkg.go.dev/github.com/haukened/fifo) [![codecov](https://codecov.io/gh/haukened/fifo/branch/main/graph/badge.svg?token=5DVM9GD3BX)](https://codecov.io/gh/haukened/fifo)

# Spool - thread safe containers.
### A collection of thread safe containers for golang that we couldn't find anywhere else

## Ring

Ring provides a fixed-length FIFO linked-list style ring structure.  `Push()` adds new elements to the ring, and `Pop()` removes the oldest element from the ring.
When the ring reaches the specified max lengh, a `Push()` operation will first remove the oldest element before adding the new element, maintaining the max length.
`Do()` accepts a function, and iterates forward performing that function on each element in the ring.

### Movement
Moving around the ring is achieved by calling `Next()`, `Prev()`, or `Move(n)` where `n` is an integer.  Positive integers `Move` forward, negative integers `Move` backwards around the ring.
All movement methods return the destination value.

### Values
Calling `Value()` on any ring element provides the value of the current cursor position within the ring.
Calling `SetValue()` on any ring element replaces the value at the current cursor position.

### Example

```
package main

import (
  "fmt"
  "bytes"
  
  "fifo"
)

func printRing(r *fifo.Ring) {
  var buf bytes.Buffer
  fmt.Fprint(&buf, "{ ")
  r.Do(func(p interface{}) {
		fmt.Fprintf(&buf, "%d ", p.(int))
	})
    fmt.Fprintf(&buf, "} Len: %d, Avail: %d", r.Len(), r.Avail())
  fmt.Println(buf.String())
}

func main() {
  max := 5
  r := fifo.NewRing(max)
  printRing(r)
  for i := 0; i < (max * 2); i++ {
    r.Push(i)
    printRing(r)
  }
  _ = r.Pop()
  printRing(r)
}

```
Results in:
```
{ } Len: 0, Avail: 5
{ 0 } Len: 1, Avail: 4
{ 0 1 } Len: 2, Avail: 3
{ 0 1 2 } Len: 3, Avail: 2
{ 0 1 2 3 } Len: 4, Avail: 1
{ 0 1 2 3 4 } Len: 5, Avail: 0
{ 1 2 3 4 5 } Len: 5, Avail: 0
{ 2 3 4 5 6 } Len: 5, Avail: 0
{ 3 4 5 6 7 } Len: 5, Avail: 0
{ 4 5 6 7 8 } Len: 5, Avail: 0
{ 5 6 7 8 9 } Len: 5, Avail: 0
{ 6 7 8 9 } Len: 4, Avail: 1
```
