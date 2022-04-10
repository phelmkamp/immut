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
	if m.RO.Len() < 1 {
		return
	}
	// No need to clone just to clear.
	// Preserve min capacity in case it matters.
	m2 := make(map[K]V, m.RO.Len())
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
	m2 := romaps.Clone(dst.RO)
	romaps.Copy(m2, src)
	dst.RO = romaps.Freeze(m2)
}

// DeleteFunc deletes any key/value pairs from m for which del returns true.
// Note: The underlying map is cloned before the write-operation is performed.
func DeleteFunc[K comparable, V any](m *Map[K, V], del func(K, V) bool) {
	// Avoid clone if pair not present.
	if !containsFunc(m.RO, del) {
		return
	}
	m2 := romaps.Clone(m.RO)
	maps.DeleteFunc(m2, del)
	m.RO = romaps.Freeze(m2)
}

func containsFunc[K comparable, V any](m romaps.Map[K, V], f func(K, V) bool) bool {
	for _, k := range romaps.Keys(m) {
		v, _ := m.Index(k)
		if f(k, v) {
			return true
		}
	}
	return false
}
