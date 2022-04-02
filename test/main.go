package main

import (
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

	// romaps
	m := romaps.Freeze(map[uint8]struct{}{})
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
}
