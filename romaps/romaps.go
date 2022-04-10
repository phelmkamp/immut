// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package romaps

import (
	"fmt"

	"golang.org/x/exp/maps"
)

// Map wraps a read-only map.
type Map[K comparable, V any] struct {
	m map[K]V
}

// Index returns the value associated with k.
// The boolean value ok is true if the value v corresponds to a key found
// in the map, false if it is a zero value because the key was not found.
func (m Map[K, V]) Index(k K) (v V, ok bool) {
	v, ok = m.m[k]
	return
}

// IsNil reports whether the underlying map is nil.
func (m Map[K, V]) IsNil() bool {
	return m.m == nil
}

// Len returns the length.
func (m Map[K, V]) Len() int {
	return len(m.m)
}

// String returns the underlying map formatted as a string.
func (m Map[K, V]) String() string {
	return fmt.Sprint(m.m)
}

// Freeze returns a read-only wrapper for the given map.
func Freeze[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V]{m: m}
}

// Clone returns a mutable copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[K comparable, V any](m Map[K, V]) map[K]V {
	return maps.Clone(m.m)
}

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst,
// the value in dst will be overwritten by the value associated
// with the key in src.
func Copy[K comparable, V any](dst map[K]V, src Map[K, V]) {
	maps.Copy(dst, src.m)
}

// Equal reports whether two maps contain the same key/value pairs.
// Values are compared using ==.
func Equal[K, V comparable](m1 Map[K, V], m2 Map[K, V]) bool {
	return maps.Equal(m1.m, m2.m)
}

// EqualFunc is like Equal, but compares values using eq.
// Keys are still compared with ==.
func EqualFunc[K comparable, V1, V2 any](m1 Map[K, V1], m2 Map[K, V2], eq func(V1, V2) bool) bool {
	return maps.EqualFunc(m1.m, m2.m, eq)
}

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[K comparable, V any](m Map[K, V]) []K {
	return maps.Keys(m.m)
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[K comparable, V any](m Map[K, V]) []V {
	return maps.Values(m.m)
}
