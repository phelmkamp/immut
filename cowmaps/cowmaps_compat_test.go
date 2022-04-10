// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cowmaps

import (
	"testing"

	"github.com/phelmkamp/immut/romaps"
)

// copied from https://cs.opensource.google/go/x/exp/+/master:slices/
// to ensure compatibility

var m1 = CopyOnWrite(map[int]int{1: 2, 2: 4, 4: 8, 8: 16})

func TestClear(t *testing.T) {
	ml := CopyOnWrite(map[int]int{1: 1, 2: 2, 3: 3})
	Clear(&ml)
	if got := ml.RO.Len(); got != 0 {
		t.Errorf("len(%v) = %d after Clear, want 0", ml, got)
	}
	if !romaps.Equal(ml.RO, romaps.Freeze((map[int]int)(nil))) {
		t.Errorf("Equal(%v, nil) = false, want true", ml)
	}
}

func TestCopy(t *testing.T) {
	mc := m1
	Copy(&mc, mc.RO)
	if !romaps.Equal(mc.RO, m1.RO) {
		t.Errorf("Copy(%v, %v) = %v, want %v", m1, m1, mc, m1)
	}
	Copy(&mc, romaps.Freeze(map[int]int{16: 32}))
	want := CopyOnWrite(map[int]int{1: 2, 2: 4, 4: 8, 8: 16, 16: 32})
	if !romaps.Equal(mc.RO, want.RO) {
		t.Errorf("Copy result = %v, want %v", mc, want)
	}
}

func TestDeleteFunc(t *testing.T) {
	mc := m1
	DeleteFunc(&mc, func(int, int) bool { return false })
	if !romaps.Equal(mc.RO, m1.RO) {
		t.Errorf("DeleteFunc(%v, true) = %v, want %v", m1, mc, m1)
	}
	DeleteFunc(&mc, func(k, v int) bool { return k > 3 })
	want := CopyOnWrite(map[int]int{1: 2, 2: 4})
	if !romaps.Equal(mc.RO, want.RO) {
		t.Errorf("DeleteFunc result = %v, want %v", mc, want)
	}
}
