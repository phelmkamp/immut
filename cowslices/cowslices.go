package cowslices

import (
	"fmt"

	"github.com/phelmkamp/immut/roslices"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Slice wraps a copy-on-write slice.
type Slice[E any] struct {
	RO roslices.Slice[E] // wraps a read-only slice
}

// String returns the underlying slice formatted as a string.
func (s Slice[E]) String() string {
	return fmt.Sprint(s.RO)
}

// CopyOnWrite returns a copy-on-write wrapper for the given slice.
func CopyOnWrite[E any](s []E) Slice[E] {
	return Slice[E]{RO: roslices.Freeze(s)}
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
// Note: The underlying slice is cloned before the write-operation is performed.
func Clip[E any](s Slice[E]) Slice[E] {
	// Avoid clone if no unused capacity to remove.
	if ro := s.RO; ro.Cap() == ro.Len() {
		return s
	}
	s2 := roslices.Clone(s.RO)
	s2 = slices.Clip(s2)
	s.RO = roslices.Freeze(s2)
	return s
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Note: The underlying slice is cloned before the write-operation is performed.
func Compact[E comparable](s Slice[E]) Slice[E] {
	// Avoid clone if already compact.
	if isCompact(s.RO) {
		return s
	}
	s2 := roslices.Clone(s.RO)
	s2 = slices.Compact(s2)
	s.RO = roslices.Freeze(s2)
	return s
}

// CompactFunc is like Compact but uses a comparison function.
// Note: The underlying slice is cloned before the write-operation is performed.
func CompactFunc[E any](s Slice[E], eq func(E, E) bool) Slice[E] {
	// Avoid clone if already compact.
	if isCompactFunc(s.RO, eq) {
		return s
	}
	s2 := roslices.Clone(s.RO)
	s2 = slices.CompactFunc(s2, eq)
	s.RO = roslices.Freeze(s2)
	return s
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete is O(len(s)-(j-i)), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Note: The underlying slice is cloned before the write-operation is performed.
func Delete[E any](s Slice[E], i, j int) Slice[E] {
	s2 := roslices.Clone(s.RO)
	s2 = slices.Delete(s2, i, j)
	s.RO = roslices.Freeze(s2)
	return s
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. Grow may modify elements of the
// slice between the length and the capacity. If n is negative or too large to
// allocate the memory, Grow panics.
// Note: The underlying slice is cloned before the write-operation is performed.
func Grow[E any](s Slice[E], n int) Slice[E] {
	// Avoid clone if capacity is already sufficient.
	ro := s.RO
	cap2 := ro.Len() + n
	if ro.Cap() >= cap2 {
		return s
	}
	// Reallocate just once with enough capacity.
	s2 := clone(ro, cap2)
	s.RO = roslices.Freeze(s2)
	return s
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// Note: The underlying slice is cloned before the write-operation is performed.
func Insert[E any](s Slice[E], i int, v ...E) Slice[E] {
	// Reallocate just once with enough capacity.
	ro := s.RO
	cap2 := ro.Cap()
	if tot := ro.Len() + len(v); tot > cap2 {
		cap2 = tot
	}
	s2 := clone(ro, cap2)
	s2 = slices.Insert(s2, i, v...)
	s.RO = roslices.Freeze(s2)
	return s
}

// Sort sorts a slice of any ordered type in ascending order.
// Note: The underlying slice is cloned before the write-operation is performed.
func Sort[E constraints.Ordered](x *Slice[E]) {
	// Avoid clone if already sorted.
	if roslices.IsSorted(x.RO) {
		return
	}
	s2 := roslices.Clone(x.RO)
	slices.Sort(s2)
	x.RO = roslices.Freeze(s2)
}

// SortFunc sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
// Note: The underlying slice is cloned before the write-operation is performed.
func SortFunc[E any](x *Slice[E], less func(a, b E) bool) {
	// Avoid clone if already sorted.
	if roslices.IsSortedFunc(x.RO, less) {
		return
	}
	s2 := roslices.Clone(x.RO)
	slices.SortFunc(s2, less)
	x.RO = roslices.Freeze(s2)
}

// SortStableFunc sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
// Note: The underlying slice is cloned before the write-operation is performed.
func SortStableFunc[E any](x *Slice[E], less func(a, b E) bool) {
	// Avoid clone if already sorted.
	if roslices.IsSortedFunc(x.RO, less) {
		return
	}
	s2 := roslices.Clone(x.RO)
	slices.SortStableFunc(s2, less)
	x.RO = roslices.Freeze(s2)
}

func isCompact[E comparable](s roslices.Slice[E]) bool {
	for i := 1; i < s.Len(); i++ {
		if s.Index(i) == s.Index(i-1) {
			return false
		}
	}
	return true
}

func isCompactFunc[E any](s roslices.Slice[E], eq func(E, E) bool) bool {
	for i := 1; i < s.Len(); i++ {
		if eq(s.Index(i), s.Index(i-1)) {
			return false
		}
	}
	return true
}

func clone[E any](ro roslices.Slice[E], cap int) []E {
	s2 := make([]E, ro.Len(), cap)
	roslices.Copy(s2, ro)
	return s2
}
