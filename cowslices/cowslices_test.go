package cowslices_test

import (
	"fmt"
	"github.com/phelmkamp/immut/cowslices"
	"github.com/phelmkamp/immut/roslices"
	"math/rand"
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

// Example_concurrent demonstrates that two concurrent goroutines
// can write to the same slice without the use of channels or locks.
func Example_concurrent() {
	rand.Seed(42)
	s := cowslices.CopyOnWrite([]int{})
	go func() {
		for {
			// append random int every 1/2 sec
			time.Sleep(time.Second / 2)
			v := rand.Intn(10)
			s = cowslices.Insert(s, s.RO.Len(), v)

		}
	}()
	go func() {
		for {
			// sort slice every 1/3 sec
			time.Sleep(time.Second / 3)
			cowslices.Sort(&s)
		}
	}()
	// run for 3 sec
	time.Sleep(3 * time.Second)
	fmt.Println(s)
	// Output: [0 3 5 7 8]
}

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
