// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package aoslices

import (
	"fmt"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Slice wraps an append-only slice.
type Slice[E any] struct {
	s []E
}

// Cap returns the capacity.
func (s Slice[E]) Cap() int {
	return cap(s.s)
}

// Index returns the i'th element.
func (s Slice[E]) Index(i int) E {
	return s.s[i]
}

// IsNil reports whether the underlying slice is nil.
func (s Slice[E]) IsNil() bool {
	return s.s == nil
}

// Len returns the length.
func (s Slice[E]) Len() int {
	return len(s.s)
}

// Slice returns s[i:].
// It panics if the indexes are out of bounds.
func (s Slice[E]) Slice(i int) Slice[E] {
	s.s = s.s[i:]
	return s
}

// String returns the underlying slice formatted as a string.
func (s Slice[E]) String() string {
	return fmt.Sprint(s.s)
}

// Make creates a new append-only Slice.
//
// The size specifies the length. The capacity of the slice is
// equal to its length. A second integer argument may be provided to
// specify a different capacity; it must be no smaller than the
// length.
func Make[E any](size ...int) Slice[E] {
	switch len(size) {
	case 1:
		return Slice[E]{s: make([]E, size[0])}
	case 2:
		return Slice[E]{s: make([]E, size[0], size[1])}
	default:
		return Slice[E]{}
	}
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i != len(s).
func Insert[E any](s Slice[E], i int, v ...E) Slice[E] {
	if i < len(s.s) {
		panic(fmt.Sprintf("aoslices.Insert: index %d prior to end of slice length %d", i, len(s.s)))
	}
	s.s = slices.Insert(s.s, i, v...)
	return s
}

// BinarySearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
func BinarySearch[E constraints.Ordered](x Slice[E], target E) (int, bool) {
	return slices.BinarySearch(x.s, target)
}

// BinarySearchFunc works like BinarySearch, but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing" is
// defined by cmp. cmp(a, b) is expected to return an integer comparing the two
// parameters: 0 if a == b, a negative number if a < b and a positive number if
// a > b.
func BinarySearchFunc[E any](x Slice[E], target E, cmp func(E, E) int) (int, bool) {
	return slices.BinarySearchFunc(x.s, target, cmp)
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[E any](s Slice[E]) Slice[E] {
	s.s = slices.Clip(s.s)
	return s
}

// Clone returns a mutable copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[E any](s Slice[E]) []E {
	return slices.Clone(s.s)
}

// Compare compares the elements of s1 and s2.
// The elements are compared sequentially, starting at index 0,
// until one element is not equal to the other.
// The result of comparing the first non-matching elements is returned.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one.
// The result is 0 if s1 == s2, -1 if s1 < s2, and +1 if s1 > s2.
// Comparisons involving floating point NaNs are ignored.
func Compare[E constraints.Ordered](s1, s2 Slice[E]) int {
	return slices.Compare(s1.s, s2.s)
}

// CompareFunc is like Compare but uses a comparison function
// on each pair of elements. The elements are compared in increasing
// index order, and the comparisons stop after the first time cmp
// returns non-zero.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func CompareFunc[E1 any, E2 any](s1 Slice[E1], s2 Slice[E2], cmp func(E1, E2) int) int {
	return slices.CompareFunc(s1.s, s2.s, cmp)
}

// Contains reports whether v is present in s.
func Contains[E comparable](s Slice[E], v E) bool {
	return slices.Contains(s.s, v)
}

// Copy copies elements from a source slice into a
// destination slice. The source and destination may overlap. Copy
// returns the number of elements copied, which will be the minimum of
// len(src) and len(dst).
func Copy[E any](dst []E, src Slice[E]) int {
	return copy(dst, src.s)
}

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func Equal[E comparable](s1, s2 Slice[E]) bool {
	return slices.Equal(s1.s, s2.s)
}

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[E1, E2 any](s1 Slice[E1], s2 Slice[E2], eq func(E1, E2) bool) bool {
	return slices.EqualFunc(s1.s, s2.s, eq)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. Grow may modify elements of the
// slice between the length and the capacity. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[E any](s Slice[E], n int) Slice[E] {
	s.s = slices.Grow(s.s, n)
	return s
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s Slice[E], v E) int {
	return slices.Index(s.s, v)
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[E any](s Slice[E], f func(E) bool) int {
	// Manually inline slices.IndexFunc so this function can be inlined by the compiler.
	for i, v := range s.s {
		if f(v) {
			return i
		}
	}
	return -1
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E constraints.Ordered](x Slice[E]) bool {
	return slices.IsSorted(x.s)
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function.
func IsSortedFunc[E any](x Slice[E], less func(a, b E) bool) bool {
	return slices.IsSortedFunc(x.s, less)
}
