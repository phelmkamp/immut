package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func Test_inline(t *testing.T) {
	want := map[string]struct{}{
		// roslices
		`inlining call to roslices\.Freeze`:                          {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].Cap`:    {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].Index`:  {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].IsNil`:  {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].Len`:    {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].Slice`:  {},
		`inlining call to roslices\.Slice\[go.shape.uint8_0].String`: {},
		`inlining call to roslices\.BinarySearch`:                    {},
		`inlining call to roslices\.BinarySearchFunc`:                {},
		`inlining call to roslices\.Clone`:                           {},
		`inlining call to roslices\.Compare`:                         {},
		`inlining call to roslices\.CompareFunc`:                     {},
		`inlining call to roslices\.Contains`:                        {},
		`inlining call to roslices\.Copy`:                            {},
		`inlining call to roslices\.Equal`:                           {},
		`inlining call to roslices\.EqualFunc`:                       {},
		`inlining call to roslices\.Index`:                           {},
		`inlining call to roslices\.IndexFunc`:                       {},
		`inlining call to roslices\.IsSorted`:                        {},
		`inlining call to roslices\.IsSortedFunc`:                    {},

		// cowslices
		`inlining call to cowslices\.CopyOnWrite`:                     {},
		`inlining call to cowslices\.Slice\[go.shape.uint8_0].String`: {},
		`inlining call to cowslices\.Delete`:                          {},
		// underlying slices calls are inlined but these cannot be
		// avoiding excess reallocation is more important anyway
		//`inlining call to cowslices\.Clip`:           {},
		//`inlining call to cowslices\.Compact`:        {},
		//`inlining call to cowslices\.CompactFunc`:    {},
		//`inlining call to cowslices\.DoAll`:          {},
		//`inlining call to cowslices\.Grow`:           {},
		//`inlining call to cowslices\.Insert`:         {},
		//`inlining call to cowslices\.Sort`:           {},
		//`inlining call to cowslices\.SortFunc`:       {},
		//`inlining call to cowslices\.SortStableFunc`: {},

		// romaps
		`inlining call to romaps\.Freeze`:                                             {},
		`inlining call to romaps\.Map\[go.shape.uint8_0,go.shape.struct {}_1].Index`:  {},
		`inlining call to romaps\.Map\[go.shape.uint8_0,go.shape.struct {}_1].IsNil`:  {},
		`inlining call to romaps\.Map\[go.shape.uint8_0,go.shape.struct {}_1].Len`:    {},
		`inlining call to romaps\.Map\[go.shape.uint8_0,go.shape.struct {}_1].String`: {},
		`inlining call to romaps\.Clone`:                                              {},
		`inlining call to romaps\.Copy`:                                               {},
		`inlining call to romaps\.Equal`:                                              {},
		`inlining call to romaps\.EqualFunc`:                                          {},
		`inlining call to romaps\.Keys`:                                               {},
		`inlining call to romaps\.Values`:                                             {},

		// cowmaps
		`inlining call to cowmaps\.CopyOnWrite`:                                        {},
		`inlining call to cowmaps\.Map\[go.shape.uint8_0,go.shape.struct {}_1].String`: {},
		`inlining call to cowmaps\.Clear`:                                              {},
		`inlining call to cowmaps\.Copy`:                                               {},
		// underlying maps calls are inlined but this cannot be
		// avoiding excess reallocation is more important anyway
		//`inlining call to cowmaps\.DeleteFunc`: {},

		// corptrs
		`inlining call to corptrs\.Freeze`:                               {},
		`inlining call to corptrs\.Ptr\[go.shape.struct {}_0]\.IsNil()`:  {},
		`inlining call to corptrs\.Ptr\[go.shape.struct {}_0]\.Clone()`:  {},
		`inlining call to corptrs\.Ptr\[go.shape.struct {}_0]\.String()`: {},
	}

	var sb strings.Builder
	cmd := exec.Command("go", "run", "-gcflags", "-m", "main.go")
	cmd.Stderr = &sb
	if err := cmd.Run(); err != nil {
		t.Fatalf("cmd.Run() failed: %v", err)
	}

	output := sb.String()
	for pattern := range want {
		matched, err := regexp.MatchString(pattern, output)
		if err != nil {
			t.Fatalf("regexp.MatchString(%s) failed: %v", pattern, err)
		}
		if matched {
			fmt.Printf("found pattern: %s\n", pattern)
			delete(want, pattern)
		}
	}

	for pattern := range want {
		t.Errorf("pattern not found: %s", pattern)
	}
}
