package roptrs_test

import (
	"fmt"

	"github.com/phelmkamp/immut/roptrs"
)

type big struct {
	a, b, c, d, e int
}

func fn1(b roptrs.Ptr[big]) {
	fn2(b)
}

func fn2(b roptrs.Ptr[big]) {
	fn3(b)
}

func fn3(b roptrs.Ptr[big]) {
	fn4(b)
}

func fn4(b roptrs.Ptr[big]) {
	fn5(b)
}

func fn5(b roptrs.Ptr[big]) {
	b2 := roptrs.Clone(b)
	b2.e++
	fmt.Println(b2)
}

func Example() {
	b := big{1, 2, 3, 4, 5}
	p := roptrs.Freeze(&b)
	fn1(p)
	fmt.Println(p)
	// Output: &{1 2 3 4 6}
	// &{1 2 3 4 5}
}
