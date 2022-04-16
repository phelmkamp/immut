package cowmaps_test

import (
	"fmt"
	"github.com/phelmkamp/immut/cowmaps"
)

func ExampleDoAll() {
	m := cowmaps.CopyOnWrite(map[string]int{"foo": 1})
	cowmaps.DoAll(&m, m.RO.Len(),
		cowmaps.DoCopy(map[string]int{"bar": 2, "baz": 3}),                              // map[foo:1 bar:2 baz:3]
		cowmaps.DoDeleteFunc[string, int](func(_ string, v int) bool { return v == 2 }), // map[foo:1 baz:3]
		cowmaps.DoSetIndex("baz", 2),                                                    // map[foo:1 baz:2]
		cowmaps.DoDelete[string, int]("foo"),                                            // map[baz:2]
	)
	fmt.Println(m)
	// Output: map[baz:2]
}
