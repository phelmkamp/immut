// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cowslices

import (
	"github.com/phelmkamp/immut/roslices"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Doer defines a method for doing an operation on a slice.
type Doer[E any] interface {
	do([]E) []E
}

// doerFunc is a func that implements Doer.
type doerFunc[E any] func([]E) []E

func (f doerFunc[E]) do(s []E) []E {
	return f(s)
}

// DoAll does all the supplied operations on the slice with minimal reallocation.
// The initial capacity of the reallocated slice is cap (or len(s) if cap is not sufficient).
// Note: The underlying slice is cloned before the write-operations are performed.
func DoAll[E any](s Slice[E], cap int, ops ...Doer[E]) Slice[E] {
	ro := s.RO
	if cap < ro.Len() {
		cap = ro.Len()
	}
	s2 := clone(ro, cap)
	for _, op := range ops {
		s2 = op.do(s2)
	}
	s.RO = roslices.Freeze(s2)
	return s
}

// DoClip returns the slices.Clip operation.
func DoClip[E any]() Doer[E] {
	return doerFunc[E](func(s []E) []E {
		return slices.Clip(s)
	})
}

// DoCompact returns the slices.Compact operation.
func DoCompact[E comparable]() Doer[E] {
	return doerFunc[E](func(s []E) []E {
		return slices.Compact(s)
	})
}

// DoCompactFunc returns the slices.CompactFunc operation.
func DoCompactFunc[E any](eq func(E, E) bool) Doer[E] {
	return doerFunc[E](func(s []E) []E {
		return slices.CompactFunc(s, eq)
	})
}

// DoDelete returns the slices.Delete operation.
func DoDelete[E any](i, j int) Doer[E] {
	return doerFunc[E](func(s []E) []E {
		return slices.Delete(s, i, j)
	})
}

// DoInsert returns the slices.Insert operation.
func DoInsert[E any](i int, v ...E) Doer[E] {
	return doerFunc[E](func(s []E) []E {
		return slices.Insert(s, i, v...)
	})
}

// DoSort returns the slices.Sort operation.
func DoSort[E constraints.Ordered]() Doer[E] {
	return doerFunc[E](func(x []E) []E {
		slices.Sort(x)
		return x
	})
}

// DoSortFunc returns the slices.SortFunc operation.
func DoSortFunc[E any](less func(a, b E) bool) Doer[E] {
	return doerFunc[E](func(x []E) []E {
		slices.SortFunc(x, less)
		return x
	})
}

// DoSortStableFunc returns the slices.SortStableFunc operation.
func DoSortStableFunc[E any](less func(a, b E) bool) Doer[E] {
	return doerFunc[E](func(x []E) []E {
		slices.SortStableFunc(x, less)
		return x
	})
}
