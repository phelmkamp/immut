package rochans

import "fmt"

// Chan wraps a read-only channel.
type Chan[T any] struct {
	ch chan T
}

func (ch Chan[T]) String() string {
	return fmt.Sprint(ch.ch)
}

// Freeze returns a read-only wrapper for the given channel.
func Freeze[T any](ch chan T) Chan[T] {
	return Chan[T]{ch: ch}
}

// Receive receives a values from the given channel.
func Receive[T any](ch Chan[T]) T {
	return <-ch.ch
}
