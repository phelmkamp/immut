package roslices

import "testing"

const N = 100_000

func fill(v, n int) []int {
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = v
	}
	return ints
}

func Benchmark_iter_index(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := Freeze(fill(1, N))
		var n int
		b.StartTimer()
		for j := 0; j < ints.Len(); j++ {
			n = j + ints.Index(j)
			n = n + 0
		}
	}
}

func Benchmark_iter_range(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := fill(1, N)
		var n int
		b.StartTimer()
		for j, v := range ints {
			n = j + v
			n = n + 0
		}
	}
}
