package romaps

import "testing"

const N = 100_000

func fill(v, n int) map[int]int {
	ints := make(map[int]int, n)
	for i := 0; i < n; i++ {
		ints[i] = v
	}
	return ints
}

func Benchmark_iter_do(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := Freeze(fill(1, N))
		var n int
		b.StartTimer()
		ints.Do(func(k, v int) bool {
			n = k + v
			n = n + 0
			return true
		})
	}
}

func Benchmark_iter_keys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := Freeze(fill(1, N))
		var n int
		b.StartTimer()
		for k := range Keys(ints) {
			v, _ := ints.Index(k)
			n = k + v
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
		for k, v := range ints {
			n = k + v
			n = n + 0
		}
	}
}
