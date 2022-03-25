# immut

[![Go Reference](https://pkg.go.dev/badge/github.com/phelmkamp/immut.svg)](https://pkg.go.dev/github.com/phelmkamp/immut)

Immutable data structures for the Go language.

## Installation

```bash
go get github.com/phelmkamp/immut
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/phelmkamp/immut/romaps"
	"github.com/phelmkamp/immut/roslices"
	_ "golang.org/x/exp/maps"
	_ "golang.org/x/exp/slices"
)

func main() {
	s := roslices.Freeze([]int{1, 2, 3})
	fmt.Println(roslices.IsSorted(s))
	//slices.Sort(s) // not allowed

	m := romaps.Freeze(map[string]int{"foo": 42, "bar": 7})
	fmt.Println(romaps.Keys(m))
	//maps.Clear(m) // not allowed
}
```
