package cowmaps

import "testing"

const N = 100_000

func fill(v, n int) map[int]int {
	ints := make(map[int]int, n)
	for i := 0; i < n; i++ {
		ints[i] = v
	}
	return ints
}

func Benchmark_containsFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := CopyOnWrite(fill(1, N))
		b.StartTimer()
		containsFunc(ints.RO, func(int, int) bool { return false })
	}
}
