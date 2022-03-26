package roptrs

import "fmt"

// Ptr wraps a read-only pointer.
type Ptr[T any] struct {
	p *T
}

func (p Ptr[T]) String() string {
	return fmt.Sprint(p.p)
}

// Freeze returns a read-only wrapper for the given pointer.
func Freeze[T any](p *T) Ptr[T] {
	return Ptr[T]{p: p}
}

// Clone creates a copy of the underlying value and returns a mutable pointer to it.
func Clone[T any](p Ptr[T]) *T {
	p2 := *p.p
	return &p2
}
