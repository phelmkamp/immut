// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package corptrs

import "fmt"

// Pointer wraps a read-only pointer.
type Pointer[T any] struct {
	p *T
}

// Clone creates a copy of the underlying value and returns a mutable pointer to it.
func (p Pointer[T]) Clone() *T {
	if p.p == nil {
		return nil
	}
	p2 := *p.p
	return &p2
}

// IsNil reports whether the underlying pointer is nil.
func (p Pointer[T]) IsNil() bool {
	return p.p == nil
}

// String returns the underlying pointer formatted as a string.
func (p Pointer[T]) String() string {
	return fmt.Sprint(p.p)
}

// Freeze returns a read-only wrapper for the given pointer.
func Freeze[T any](p *T) Pointer[T] {
	return Pointer[T]{p: p}
}
