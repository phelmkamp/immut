# immut

[![Go Reference](https://pkg.go.dev/badge/github.com/phelmkamp/immut.svg)](https://pkg.go.dev/github.com/phelmkamp/immut)
[![Go Report Card](https://goreportcard.com/badge/github.com/phelmkamp/immut)](https://goreportcard.com/report/github.com/phelmkamp/immut)
[![codecov](https://codecov.io/gh/phelmkamp/immut/branch/main/graph/badge.svg?token=79CVDP412S)](https://codecov.io/gh/phelmkamp/immut)

In Go, immutability is limited to primitive types and structs (via the `const` keyword and pass-by-value respectively).
This module provides read-only slices, maps, and pointers via the low/zero-cost abstractions `roslices.Slice`, `romaps.Map`, and `corptrs.Ptr`.
Go 1.18+ parameterized types (AKA generics) are used to support any underlying type.

In addition, the `cowslices` and `cowmaps` packages provide copy-on-write semantics. The mutating functions clone the underlying value before the write-operation is performed

The `*slices` and `*maps` packages are drop-in replacements for the standard [slices](https://pkg.go.dev/golang.org/x/exp/slices) and [maps](https://pkg.go.dev/golang.org/x/exp/maps) packages.

Note: This project is not intended to replace all uses of slices, maps, and pointers with unnecessary boxed-types. It is specifically for restricting write-access to values of these types in cases where it is desirable to do so.
For example:
 * Ensure return values cannot be modified (previously required a defensive copy)
 * Guarantee to callers that function arguments will not change
 * Safely access shared state from multiple goroutines (e.g. values in `context.Context`)
 * Prevent mutation of variables/fields after initialization
 * Pass pointers to large structs and avoid excess copying without the risk of modification


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
	"github.com/phelmkamp/immut/romaps"
	"github.com/phelmkamp/immut/corptrs"
	"github.com/phelmkamp/immut/roslices"
	//"golang.org/x/exp/maps"
	//"golang.org/x/exp/slices"
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

	// read-only pointers
	type big struct {
		a, b, c, d, e int
	}
	b := big{1, 2, 3, 4, 5}
	p := corptrs.Freeze(&b)
	p2 := p.Clone()
	p2.a = 42
	fmt.Println(p, p2)
	//p.a = 42 // not allowed
}
```

## Releases

This module strives to maintain compatibility with packages currently in [exp](https://pkg.go.dev/golang.org/x/exp).
As such, it will remain untagged until the corresponding packages are tagged.
Then, it will remain at v0 until the corresponding packages achieve v1 status.

## Performance

This project aims to be a low/zero-cost abstraction of Go's standard types.
The compiler can inline almost all read-only function calls as verified by `Test_inline` in the `test` submodule.
 * Exceptions: `roslices.IndexFunc`

The copy-on-write functions avoid unnecessary reallocation wherever possible.
As such, most of the copy-on-write functions cannot be inlined by the compiler but that is a conscious tradeoff.
The `cowslices.DoAll` function is provided to support multiple write-operations with minimal reallocation.
