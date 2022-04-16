// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cowmaps

import (
	"fmt"

	"github.com/phelmkamp/immut/romaps"
	"golang.org/x/exp/maps"
)

// Map wraps a copy-on-write map.
type Map[K comparable, V any] struct {
	RO romaps.Map[K, V] // wraps a read-only map
}

// Delete deletes the element with the specified key from the map.
// If there is no such element, delete is a no-op.
// Note: The underlying map is reallocated before the write-operation is performed.
func (m *Map[K, V]) Delete(k K) (v V, ok bool) {
	// Avoid reallocation if key not present.
	ro := m.RO
	if v, ok = ro.Index(k); !ok {
		return
	}
	m2 := romaps.Clone(ro)
	delete(m2, k)
	m.RO = romaps.Freeze(m2)
	return
}

// SetIndex sets the element associated with k to v.
// Note: The underlying map is reallocated before the write-operation is performed.
func (m *Map[K, V]) SetIndex(k K, v V) {
	ro := m.RO
	m2 := clone(ro, ro.Len()+1)
	m2[k] = v
	m.RO = romaps.Freeze(m2)
}

// String returns the underlying map formatted as a string.
func (m Map[K, V]) String() string {
	return fmt.Sprint(m.RO)
}

// CopyOnWrite returns a copy-on-write wrapper for the given map.
func CopyOnWrite[K comparable, V any](m map[K]V) Map[K, V] {
	return Map[K, V]{RO: romaps.Freeze(m)}
}

// Clear removes all entries from m, leaving it empty.
// Note: The underlying map is reallocated before the write-operation is performed.
func Clear[K comparable, V any](m *Map[K, V]) {
	// Avoid reallocation if m is empty.
	ro := m.RO
	if ro.Len() < 1 {
		return
	}
	// No need to clone just to clear.
	m2 := make(map[K]V)
	m.RO = romaps.Freeze(m2)
}

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst,
// the value in dst will be overwritten by the value associated
// with the key in src.
// Note: The underlying map is cloned before the write-operation is performed.
func Copy[K comparable, V any](dst *Map[K, V], src romaps.Map[K, V]) {
	// Avoid clone if src is empty.
	if src.Len() < 1 {
		return
	}
	ro := dst.RO
	m2 := clone(ro, ro.Len()+src.Len()) // Ensure no additional allocation.
	romaps.Copy(m2, src)
	dst.RO = romaps.Freeze(m2)
}

// DeleteFunc deletes any key/value pairs from m for which del returns true.
// Note: The underlying map is cloned before the write-operation is performed.
func DeleteFunc[K comparable, V any](m *Map[K, V], del func(K, V) bool) {
	// Avoid clone if pair not present.
	ro := m.RO
	if !containsFunc(ro, del) {
		return
	}
	m2 := romaps.Clone(ro)
	maps.DeleteFunc(m2, del)
	m.RO = romaps.Freeze(m2)
}

func clone[K comparable, V any](m romaps.Map[K, V], cap int) map[K]V {
	m2 := make(map[K]V, cap)
	romaps.Copy(m2, m)
	return m2
}

func containsFunc[K comparable, V any](m romaps.Map[K, V], f func(K, V) bool) bool {
	var found bool
	m.Do(func(k K, v V) bool {
		found = f(k, v)
		return !found
	})
	return found
}
