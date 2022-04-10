// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cowslices_test

import (
	"fmt"
	"github.com/phelmkamp/immut/cowslices"
	"reflect"
	"testing"
)

func ExampleDoAll() {
	s := cowslices.CopyOnWrite([]int{1, 2, 2, 3})
	s = cowslices.DoAll(s, s.RO.Len(),
		cowslices.DoInsert[int](1, 3), // [1 3 2 2 3]
		cowslices.DoSort[int](),       // [1 2 2 3 3]
		cowslices.DoCompact[int](),    // [1 2 3]
		cowslices.DoDelete[int](1, 2), // [1 3]
		cowslices.DoClip[int](),       // [1 3]
	)
	fmt.Println(s)
	// Output: [1 3]
}

func TestDoAll(t *testing.T) {
	s1 := cowslices.CopyOnWrite([]int{1, 2, 2, 3})
	got := cowslices.DoAll(s1, s1.RO.Len(),
		cowslices.DoInsert[int](1, 3), // [1 3 2 2 3]
		cowslices.DoSort[int](),       // [1 2 2 3 3]
		cowslices.DoCompact[int](),    // [1 2 3]
		cowslices.DoDelete[int](1, 2), // [1 3]
		cowslices.DoClip[int](),       // [1 3]
	)
	if want := cowslices.CopyOnWrite([]int{1, 3}); !reflect.DeepEqual(got, want) {
		t.Errorf("All() = %v, want %v", got, want)
	}
	if want := 2; got.RO.Cap() != want {
		t.Errorf("All().RO.Cap() = %v, want %v", got.RO.Cap(), want)
	}
}

func TestDoAll_func(t *testing.T) {
	s1 := cowslices.CopyOnWrite([][]int{{1}, {1, 2}, {1}})
	got := cowslices.DoAll(s1, s1.RO.Len(),
		cowslices.DoSortFunc[[]int](func(a, b []int) bool { return len(a) < len(b) }),            // [[1] [1] [1 2]]
		cowslices.DoCompactFunc[[]int](func(a, b []int) bool { return reflect.DeepEqual(a, b) }), // [[1] [1 2]]
	)
	if want := cowslices.CopyOnWrite([][]int{{1}, {1, 2}}); !reflect.DeepEqual(got, want) {
		t.Errorf("All() = %v, want %v", got, want)
	}
}
