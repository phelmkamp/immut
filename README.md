# immut

[![Go Reference](https://pkg.go.dev/badge/github.com/phelmkamp/immut.svg)](https://pkg.go.dev/github.com/phelmkamp/immut)
[![Go Report Card](https://goreportcard.com/badge/github.com/phelmkamp/immut)](https://goreportcard.com/report/github.com/phelmkamp/immut)
[![codecov](https://codecov.io/gh/phelmkamp/immut/branch/main/graph/badge.svg?token=79CVDP412S)](https://codecov.io/gh/phelmkamp/immut)

Immutable data structures for the Go language.

Utilizes generics instead of reflection.
All types are "reference" types; the underlying values are not reallocated unless clearly stated (e.g. clone, copy).

Designed to be a drop-in replacement for [slices](https://pkg.go.dev/golang.org/x/exp/slices) and [maps](https://pkg.go.dev/golang.org/x/exp/maps) with most functions delegating to those packages.

## Installation

```bash
go get github.com/phelmkamp/immut
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/phelmkamp/immut/cowmaps"
	"github.com/phelmkamp/immut/cowslices"
	"github.com/phelmkamp/immut/rochans"
	"github.com/phelmkamp/immut/romaps"
	"github.com/phelmkamp/immut/roptrs"
	"github.com/phelmkamp/immut/roslices"
	_ "golang.org/x/exp/maps"
	_ "golang.org/x/exp/slices"
)

func main() {
	// read-only slices
	s := roslices.Freeze([]int{1, 2, 3})
	fmt.Println(roslices.IsSorted(s))
	//slices.Sort(s) // not allowed

	// copy-on-write slices
	s2 := cowslices.CopyOnWrite([]int{2, 1, 3})
	cowslices.Sort(&s2)
	fmt.Println(s2)

	// read-only maps
	m := romaps.Freeze(map[string]int{"foo": 42, "bar": 7})
	fmt.Println(romaps.Keys(m))
	//maps.Clear(m) // not allowed

	// copy-on-write maps
	m2 := cowmaps.CopyOnWrite(map[string]int{"foo": 42, "bar": 7})
	cowmaps.DeleteFunc(&m2, func(k string, v int) bool { return k == "foo" })
	fmt.Println(m2)

	// read-only channels
	ch := make(chan int)
	roch := rochans.Freeze(ch)
	go func() {
		ch <- 42
	}()
	fmt.Println(roch.Recv())
	//roch <- 7 // not allowed

	// read-only pointers
	type big struct {
		a, b, c, d, e int
	}
	b := big{1, 2, 3, 4, 5}
	p := roptrs.Freeze(&b)
	p2 := p.Clone()
	p2.a = 42
	fmt.Println(p, p2)
	//p.a = 42 // not allowed
}
```

## Releases

This module strives to maintain compatibility with packages currently in [exp](https://pkg.go.dev/golang.org/x/exp).
As such, it will remain untagged until the corresponding packages are tagged.
Then, it will remain at v0.x.x until the corresponding packages achieve v1.x.x status.