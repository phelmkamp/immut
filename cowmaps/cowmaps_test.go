package cowmaps_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/phelmkamp/immut/cowmaps"
	"github.com/phelmkamp/immut/romaps"
)

func Example() {
	m1 := map[string]int{"fiz": 3}
	m2 := cowmaps.CopyOnWrite(m1)
	cowmaps.Copy(&m2, romaps.Freeze(map[string]int{"foo": 42, "bar": 7}))
	fmt.Println(m2)
	fmt.Println(m1)
	// Output: map[bar:7 fiz:3 foo:42]
	// map[fiz:3]
}

func makeMap(n int) map[string]*int {
	rand.Seed(42)
	m := make(map[string]*int, n)
	for i := 0; i < n; i++ {
		v := rand.Intn(n)
		m[strconv.Itoa(v)] = &v
	}
	return m
}

// Example_concurrent demonstrates that two concurrent goroutines
// can access the same map without the use of channels or locks.
func Example_concurrent() {
	m := cowmaps.CopyOnWrite(makeMap(5_000))
	go func() {
		for {
			// delete 1 pair after slight delay
			time.Sleep(1 * time.Millisecond)
			kDel := romaps.Keys(m.RO)[0]
			cowmaps.DeleteFunc(&m, func(k string, v *int) bool {
				return k == kDel
			})
		}
	}()
	go func() {
		for {
			// read all pairs constantly
			// without COW panic is possible
			// but ro is guaranteed not to change
			ro := m.RO
			for _, k := range romaps.Keys(ro) {
				v, _ := ro.Index(k)
				_ = fmt.Sprint(*v)
			}
		}
	}()
	// run for 1 sec
	time.Sleep(1 * time.Second)
	// Output:
}

// Example_concurrent_mutable is Example_concurrent written with regular maps.
// Uncomment and run to observe panic when reading map.
//func Example_concurrent_mutable() {
//	m := makeMap(5_000)
//	go func() {
//		for {
//			// delete 1 pair after slight delay
//			time.Sleep(1 * time.Millisecond)
//			kDel := maps.Keys(m)[0]
//			maps.DeleteFunc(m, func(k string, v *int) bool {
//				return k == kDel
//			})
//		}
//	}()
//	go func() {
//		for {
//			// read all pairs constantly
//			// without COW panic is possible
//			for _, k := range maps.Keys(m) {
//				v, _ := m[k]
//				_ = fmt.Sprint(*v)
//			}
//		}
//	}()
//	// run for 1 sec
//	time.Sleep(1 * time.Second)
//	// Output:
//}

func TestMap_String(t *testing.T) {
	type fields struct {
		m map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				m: nil,
			},
			want: "map[]",
		},
		{
			name: "empty",
			fields: fields{
				m: make(map[string]int),
			},
			want: "map[]",
		},
		{
			name: "one",
			fields: fields{
				m: map[string]int{"foo": 42},
			},
			want: "map[foo:42]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := cowmaps.CopyOnWrite(tt.fields.m)
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
