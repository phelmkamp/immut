// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cowslices

import (
	"golang.org/x/exp/slices"
	"testing"
)

func fill(v, n int) []int {
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = v
	}
	return ints
}

func BenchmarkCompact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := CopyOnWrite(fill(1, N))
		b.StartTimer()
		Compact(ints)
	}
}

func BenchmarkSlicesCompact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := fill(1, N)
		b.StartTimer()
		slices.Compact(ints)
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := CopyOnWrite(fill(1, N))
		b.StartTimer()
		Insert(ints, ints.RO.Len()-1, ints.RO.Len())
	}
}

func BenchmarkSlicesInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := fill(1, N)
		b.StartTimer()
		slices.Insert(ints, len(ints)-1, len(ints))
	}
}
