package roslices

import (
	"fmt"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Slice wraps a read-only slice.
type Slice[E any] struct {
	s []E
}

func (s Slice[E]) Cap() int {
	return cap(s.s)
}

func (s Slice[E]) Index(i int) E {
	return s.s[i]
}

func (s Slice[E]) IsNil() bool {
	return s.s == nil
}

func (s Slice[E]) Len() int {
	return len(s.s)
}

func (s Slice[E]) String() string {
	return fmt.Sprint(s.s)
}

// Freeze returns a read-only wrapper for the given slice.
func Freeze[E any](s []E) Slice[E] {
	return Slice[E]{s: s}
}

// BinarySearch searches for target in a sorted slice and returns the smallest
// index at which target is found. If the target is not found, the index at
// which it could be inserted into the slice is returned; therefore, if the
// intention is to find target itself a separate check for equality with the
// element at the returned index is required.
func BinarySearch[E constraints.Ordered](x Slice[E], target E) int {
	return slices.BinarySearch(x.s, target)
}

// BinarySearchFunc uses binary search to find and return the smallest index i
// in [0, n) at which ok(i) is true, assuming that on the range [0, n),
// ok(i) == true implies ok(i+1) == true. That is, BinarySearchFunc requires
// that ok is false for some (possibly empty) prefix of the input range [0, n)
// and then true for the (possibly empty) remainder; BinarySearchFunc returns
// the first true index. If there is no such index, BinarySearchFunc returns n.
// (Note that the "not found" return value is not -1 as in, for instance,
// strings.Index.) Search calls ok(i) only for i in the range [0, n).
func BinarySearchFunc[E any](s Slice[E], ok func(E) bool) int {
	return slices.BinarySearchFunc(s.s, ok)
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

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s Slice[E], v E) int {
	return slices.Index(s.s, v)
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[E any](s Slice[E], f func(E) bool) int {
	return slices.IndexFunc(s.s, f)
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
