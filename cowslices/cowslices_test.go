package cowslices_test

import (
	"fmt"

	"github.com/phelmkamp/immut/cowslices"
	"github.com/phelmkamp/immut/roslices"
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
