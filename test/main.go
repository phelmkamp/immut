// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/phelmkamp/immut/corptrs"
	"github.com/phelmkamp/immut/cowmaps"
	"github.com/phelmkamp/immut/cowslices"
	"github.com/phelmkamp/immut/romaps"
	"github.com/phelmkamp/immut/roslices"
)

func main() {
	// roslices
	s := roslices.Freeze([]uint8{0})
	s.Cap()
	s.Index(0)
	s.IsNil()
	s.Len()
	s.Slice(0, 1)
	s.String()
	roslices.BinarySearch(s, 0)
	roslices.BinarySearchFunc(s, 0, func(uint8, uint8) int { return 0 })
	roslices.Clone(s)
	roslices.Compare(s, s)
	roslices.CompareFunc(s, s, func(uint8, uint8) int { return 0 })
	roslices.Contains(s, 0)
	roslices.Copy(nil, s)
	roslices.Equal(s, s)
	roslices.EqualFunc(s, s, func(uint8, uint8) bool { return true })
	roslices.Index(s, 0)
	roslices.IndexFunc(s, func(uint8) bool { return true })
	roslices.IsSorted(s)
	roslices.IsSortedFunc(s, func(uint8, uint8) bool { return true })

	// cowslices
	cs := cowslices.CopyOnWrite([]uint8{0})
	cs.String()
	cowslices.Clip(cs)
	cowslices.Compact(cs)
	cowslices.CompactFunc(cs, func(uint8, uint8) bool { return true })
	cowslices.Delete(cs, 0, 1)
	cowslices.DoAll(cs, -1, cowslices.DoSort[uint8]())
	cowslices.Grow(cs, 1)
	cowslices.Insert(cs, 0, 0)
	cowslices.Sort(&cs)
	cowslices.SortFunc(&cs, func(uint8, uint8) bool { return true })
	cowslices.SortStableFunc(&cs, func(uint8, uint8) bool { return true })

	// romaps
	m := romaps.Freeze(map[uint8]struct{}{})
	m.Do(func(uint8, struct{}) bool { return false })
	m.Index(0)
	m.IsNil()
	m.Len()
	m.String()
	romaps.Clone(m)
	romaps.Copy(nil, m)
	romaps.Equal(m, m)
	romaps.EqualFunc(m, m, func(struct{}, struct{}) bool { return true })
	romaps.Keys(m)
	romaps.Values(m)

	// cowmaps
	cm := cowmaps.CopyOnWrite(map[uint8]struct{}{})
	cm.String()
	cowmaps.Clear(&cm)
	cowmaps.Copy(&cm, m)
	cowmaps.DeleteFunc(&cm, func(uint8, struct{}) bool { return false })

	// corptrs
	p := corptrs.Freeze(&struct{}{})
	p.IsNil()
	p.Clone()
	p.String()
}
