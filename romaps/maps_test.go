package romaps_test

import (
	"fmt"

	"github.com/phelmkamp/immut/romaps"
)

func Example() {
	m1 := romaps.Freeze(map[string]int{"foo": 42, "bar": 7})
	m2 := map[string]int{"fiz": 3}
	fmt.Println(len(romaps.Keys(m1)))
	romaps.Copy(m2, m1)
	fmt.Println(m2)
	// Output: 2
	// map[bar:7 fiz:3 foo:42]
}
