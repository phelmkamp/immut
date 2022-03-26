package rochans

import "fmt"

// Chan wraps a read-only channel.
type Chan[T any] struct {
	ch chan T
}

func (ch Chan[T]) Cap() int {
	return cap(ch.ch)
}

func (ch Chan[T]) IsNil() bool {
	return ch.ch == nil
}

func (ch Chan[T]) Len() int {
	return len(ch.ch)
}

func (ch Chan[T]) String() string {
	return fmt.Sprint(ch.ch)
}

// Freeze returns a read-only wrapper for the given channel.
func Freeze[T any](ch chan T) Chan[T] {
	return Chan[T]{ch: ch}
}

// Recv receives a value from the given channel.
func (ch Chan[T]) Recv() (T, bool) {
	v, ok := <-ch.ch
	return v, ok
}
