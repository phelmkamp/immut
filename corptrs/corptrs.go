package corptrs

import "fmt"

// Ptr wraps a read-only pointer.
type Ptr[T any] struct {
	p *T
}

// Clone creates a copy of the underlying value and returns a mutable pointer to it.
func (p Ptr[T]) Clone() *T {
	if p.p == nil {
		return nil
	}
	p2 := *p.p
	return &p2
}

// IsNil reports whether the underlying pointer is nil.
func (p Ptr[T]) IsNil() bool {
	return p.p == nil
}

// String returns the underlying pointer formatted as a string.
func (p Ptr[T]) String() string {
	return fmt.Sprint(p.p)
}

// Freeze returns a read-only wrapper for the given pointer.
func Freeze[T any](p *T) Ptr[T] {
	return Ptr[T]{p: p}
}
