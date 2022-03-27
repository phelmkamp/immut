package rochans

import "fmt"

// Chan wraps a read-only channel.
type Chan[T any] struct {
	ch chan T
}

// Cap returns the capacity.
func (ch Chan[T]) Cap() int {
	return cap(ch.ch)
}

// IsNil reports whether the underlying channel is nil.
func (ch Chan[T]) IsNil() bool {
	return ch.ch == nil
}

// Len returns the length.
func (ch Chan[T]) Len() int {
	return len(ch.ch)
}

// Recv receives and returns a value from the channel.
// The receive blocks until a value is ready.
// The boolean value ok is true if the value v corresponds to a send
// on the channel, false if it is a zero value received because the channel is closed.
func (ch Chan[T]) Recv() (v T, ok bool) {
	v, ok = <-ch.ch
	return
}

// String returns the underlying channel formatted as a string.
func (ch Chan[T]) String() string {
	return fmt.Sprint(ch.ch)
}

// Freeze returns a read-only wrapper for the given channel.
func Freeze[T any](ch chan T) Chan[T] {
	return Chan[T]{ch: ch}
}
