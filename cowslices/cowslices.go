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
	s2 := roslices.Clone(s.RO)
	s2 = slices.Clip(s2)
	s.RO = roslices.Freeze(s2)
	return s
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Note: The underlying slice is cloned before the write-operation is performed.
func Compact[E comparable](s Slice[E]) Slice[E] {
	s2 := roslices.Clone(s.RO)
	s2 = slices.Compact(s2)
	s.RO = roslices.Freeze(s2)
	return s
}

// CompactFunc is like Compact but uses a comparison function.
// Note: The underlying slice is cloned before the write-operation is performed.
func CompactFunc[E any](s Slice[E], eq func(E, E) bool) Slice[E] {
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
	s2 := roslices.Clone(s.RO)
	s2 = slices.Grow(s2, n)
	s.RO = roslices.Freeze(s2)
	return s
}

// Sort sorts a slice of any ordered type in ascending order.
// Note: The underlying slice is cloned before the write-operation is performed.
func Sort[E constraints.Ordered](x *Slice[E]) {
	s2 := roslices.Clone(x.RO)
	slices.Sort(s2)
	x.RO = roslices.Freeze(s2)
}

// SortFunc sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
// Note: The underlying slice is cloned before the write-operation is performed.
func SortFunc[E any](x *Slice[E], less func(a, b E) bool) {
	s2 := roslices.Clone(x.RO)
	slices.SortFunc(s2, less)
	x.RO = roslices.Freeze(s2)
}

// SortStableFunc sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
// Note: The underlying slice is cloned before the write-operation is performed.
func SortStableFunc[E any](x *Slice[E], less func(a, b E) bool) {
	s2 := roslices.Clone(x.RO)
	slices.SortStableFunc(s2, less)
	x.RO = roslices.Freeze(s2)
}
