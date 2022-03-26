package rochans_test

import (
	"fmt"

	"github.com/phelmkamp/immut/rochans"
)

func Example() {
	ch := make(chan int)
	roch := rochans.Freeze(ch)
	go func() {
		ch <- 42
	}()
	fmt.Println(rochans.Receive(roch))
	// Output: 42
}
