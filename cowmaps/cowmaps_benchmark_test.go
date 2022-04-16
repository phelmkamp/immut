package cowmaps

import (
	"math"
	"sync"
	"testing"
)

const N = 100_000

func fill(v, n int) map[int]int {
	ints := make(map[int]int, n)
	for i := 0; i < n; i++ {
		ints[i] = v
	}
	return ints
}

func Benchmark_containsFunc(b *testing.B) {
	ints := CopyOnWrite(fill(1, N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsFunc(ints.RO, func(int, int) bool { return false })
	}
}

func Benchmark(b *testing.B) {
	m := CopyOnWrite(make(map[int]int))
	// seed one entry for stability
	m.SetIndex(1, 1)
	ratio := int(math.Ceil(0.05 * float64(b.N))) // 1/5 miss
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(&m, b.N, ratio, nil)
	}
}

func Benchmark_sync(b *testing.B) {
	m := &sync.Map{}
	// seed one entry for stability
	m.Store(1, 1)
	ratio := int(math.Ceil(0.05 * float64(b.N))) // 1/5 miss
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runSync(m, b.N, ratio, nil)
	}
}

func run(m *Map[int, int], n, ratio int, f func(int, bool)) {
	mu := &sync.Mutex{}
	for i := 0; i < n; i++ {
		k := 1
		if i%ratio == 0 {
			// miss
			k = i
		}
		v, ok := m.RO.Index(k)
		if !ok {
			v = i
			go func() {
				mu.Lock()
				defer mu.Unlock()
				m.SetIndex(i, i)
			}()
		}
		if f != nil {
			f(v, ok)
		}
	}
}

func runSync(m *sync.Map, n, ratio int, f func(any, bool)) {
	for i := 0; i < n; i++ {
		k := 1
		if i%ratio == 0 {
			// miss
			k = i
		}
		v, ok := m.Load(k)
		if !ok {
			v = i
			go func() { m.Store(i, i) }()
		}
		if f != nil {
			f(v, ok)
		}
	}
}

//func Test(t *testing.T) {
//	m := CopyOnWrite(make(map[int]int))
//	// seed one entry for stability
//	m.SetIndex(1, 1)
//	n := 10_000
//	ratio := int(math.Ceil(0.05 * float64(n))) // 1/5 miss
//	f := func(v int, ok bool) {
//		//fmt.Println(v, ok)
//		if (ok && v%ratio == 0) || (!ok && v%ratio != 0) {
//			t.Errorf("unexpected cache result: %v %v", v, ok)
//		}
//	}
//	run(&m, n, ratio, f)
//	fmt.Println(m)
//}
//
//func Test_sync(t *testing.T) {
//	m := &sync.Map{}
//	// seed one entry for stability
//	m.Store(1, 1)
//	n := 10_000
//	ratio := int(math.Ceil(0.05 * float64(n))) // 1/5 miss
//	f := func(v any, ok bool) {
//		//fmt.Println(v, ok)
//		if (ok && v.(int)%ratio == 0) || (!ok && v.(int)%ratio != 0) {
//			t.Errorf("unexpected cache result: %v %v", v, ok)
//		}
//	}
//	runSync(m, n, ratio, f)
//	fmt.Println(m)
//}
