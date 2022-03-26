package cowmaps_test

import (
	"fmt"

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
