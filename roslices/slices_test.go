package roslices_test

import (
	"fmt"
	"testing"

	"github.com/phelmkamp/immut/roslices"
	"golang.org/x/exp/slices"
)

func Example() {
	ints1 := roslices.Freeze([]int{2, 1, 3})
	if !roslices.IsSorted(ints1) {
		// must clone to sort
		ints2 := roslices.Clone(ints1)
		slices.Sort(ints2)
		fmt.Println(ints1, ints2)
	}
	// Output: [2 1 3] [1 2 3]
}

func TestBinarySearchFunc(t *testing.T) {
	ints := roslices.Freeze([]int{1, 2, 3})
	got := roslices.BinarySearchFunc(ints, func(i int) bool {
		return i == 3
	})
	if got != 2 {
		t.Errorf("BinarySearchFunc() = %v, want %v", got, 2)
	}

	strings := roslices.Freeze([]string{"1", "2", "3"})
	got = roslices.BinarySearchFunc(strings, func(s string) bool {
		return s == "3"
	})
	if got != 2 {
		t.Errorf("BinarySearchFunc() = %v, want %v", got, 2)
	}
}

func TestBinarySearch(t *testing.T) {
	ints := roslices.Freeze([]int{1, 2, 3})
	got := roslices.BinarySearch(ints, 3)
	if got != 2 {
		t.Errorf("BinarySearch() = %v, want %v", got, 2)
	}

	//ifcs := roslices.Freeze([]interface{}{1, 2, 3})
	//got := roslices.BinarySearch(ifcs, 3)
	//if got != 2 {
	//	t.Errorf("BinarySearch() = %v, want %v", got, 2)
	//}
}
