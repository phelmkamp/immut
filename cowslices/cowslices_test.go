package cowslices_test

import (
	"fmt"
	"github.com/phelmkamp/immut/cowslices"
	"github.com/phelmkamp/immut/roslices"
	"testing"
	"time"
)

func Example() {
	ints1 := []int{2, 1, 3}
	ints2 := cowslices.CopyOnWrite(ints1)
	if !roslices.IsSorted(ints2.RO) {
		cowslices.Sort(&ints2) // ints1 is not affected
		fmt.Println(ints2)
	}
	fmt.Println(ints1)
	// Output: [1 2 3]
	// [2 1 3]
}

func makeInts(n int) []int {
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = i
	}
	return ints
}

// Example_concurrent demonstrates that two concurrent goroutines
// can access the same slice without the use of channels or locks.
func Example_concurrent() {
	s := cowslices.CopyOnWrite(makeInts(10_000))
	go func() {
		for {
			// delete element 1 after slight delay
			time.Sleep(1 * time.Millisecond)
			s = cowslices.Delete(s, 1, 2)
		}
	}()
	go func() {
		for {
			// read last element constantly
			// without COW index out-of-bounds is possible
			// but ro is guaranteed not to change
			ro := s.RO
			_ = fmt.Sprint(ro.Index(ro.Len() - 1))
		}
	}()
	// run for 5 sec
	time.Sleep(5 * time.Second)
	fmt.Println(s.RO.Index(0))
	// Output: 0
}

// Example_concurrent_mutable is Example_concurrent written with regular slices.
// Uncomment and run to observe panic when reading slice.
//func Example_concurrent_mutable() {
//	s := makeInts(10_000)
//	go func() {
//		for {
//			// delete element 1 after slight delay
//			time.Sleep(1 * time.Millisecond)
//			s = append(s[:1], s[2:]...)
//		}
//	}()
//	go func() {
//		for {
//			// read last element constantly
//			// without COW index out-of-bounds is possible
//			_ = fmt.Sprint(s[len(s)-1])
//		}
//	}()
//	// run for 5 sec
//	time.Sleep(5 * time.Second)
//	fmt.Println(s[0])
//	// Output: 0
//}

func TestSlice_String(t *testing.T) {
	type fields struct {
		s []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			want: "[]",
		},
		{
			name: "empty",
			fields: fields{
				s: make([]int, 0),
			},
			want: "[]",
		},
		{
			name: "one",
			fields: fields{
				s: []int{0},
			},
			want: "[0]",
		},
		{
			name: "two",
			fields: fields{
				s: []int{0, 1},
			},
			want: "[0 1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := cowslices.CopyOnWrite(tt.fields.s)
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
