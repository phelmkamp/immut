// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cowslices

import (
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/phelmkamp/immut/roslices"
)

// copied from https://cs.opensource.google/go/x/exp/+/master:slices/
// to ensure compatibility

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3}
var float64sWithNaNs = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8}
var strs = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestSortIntSlice(t *testing.T) {
	data := CopyOnWrite(ints[:])
	Sort(&data)
	if !roslices.IsSorted(data.RO) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFuncIntSlice(t *testing.T) {
	data := CopyOnWrite(ints[:])
	SortFunc(&data, func(a, b int) bool { return a < b })
	if !roslices.IsSorted(data.RO) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := CopyOnWrite(float64s[:])
	Sort(&data)
	if !roslices.IsSorted(data.RO) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64SliceWithNaNs(t *testing.T) {
	data := float64sWithNaNs[:]
	cowData := CopyOnWrite(data)
	input := make([]float64, len(float64sWithNaNs))
	for i := range input {
		input[i] = float64sWithNaNs[i]
	}
	// Make sure Sort doesn't panic when the slice contains NaNs.
	Sort(&cowData)
	data = roslices.Clone(cowData.RO)
	// Check whether the result is a permutation of the input.
	sort.Float64s(data)
	sort.Float64s(input)
	for i, v := range input {
		if data[i] != v && !(math.IsNaN(data[i]) && math.IsNaN(v)) {
			t.Fatalf("the result is not a permutation of the input\ngot %v\nwant %v", data, input)
		}
	}
}

func TestSortStringSlice(t *testing.T) {
	data := CopyOnWrite(strs[:])
	Sort(&data)
	if !roslices.IsSorted(data.RO) {
		t.Errorf("sorted %v", strs)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	d := make([]int, n)
	for i := 0; i < len(d); i++ {
		d[i] = rand.Intn(100)
	}
	data := CopyOnWrite(d)
	if roslices.IsSorted(data.RO) {
		t.Fatalf("terrible rand.rand")
	}
	Sort(&data)
	if !roslices.IsSorted(data.RO) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

type intPair struct {
	a, b int
}

type intPairs []intPair

// Pairs compare on a only.
func intPairLess(x, y intPair) bool {
	return x.a < y.a
}

// Record initial order in B.
func (d intPairs) initB() {
	for i := range d {
		d[i].b = i
	}
}

// InOrder checks if a-equal elements were not reordered.
func (d intPairs) inOrder() bool {
	lastA, lastB := -1, 0
	for i := 0; i < len(d); i++ {
		if lastA != d[i].a {
			lastA = d[i].a
			lastB = d[i].b
			continue
		}
		if d[i].b <= lastB {
			return false
		}
		lastB = d[i].b
	}
	return true
}

func TestStability(t *testing.T) {
	n, m := 100000, 1000
	if testing.Short() {
		n, m = 1000, 100
	}
	d := make(intPairs, n)

	// random distribution
	for i := 0; i < len(d); i++ {
		d[i].a = rand.Intn(m)
	}
	data := CopyOnWrite(d)
	if roslices.IsSortedFunc(data.RO, intPairLess) {
		t.Fatalf("terrible rand.rand")
	}
	d.initB()
	SortStableFunc(&data, intPairLess)
	if !roslices.IsSortedFunc(data.RO, intPairLess) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	d = roslices.Clone(data.RO)
	if !d.inOrder() {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}

	// already sorted
	d.initB()
	SortStableFunc(&data, intPairLess)
	if !roslices.IsSortedFunc(data.RO, intPairLess) {
		t.Errorf("Stable shuffled sorted %d ints (order)", n)
	}
	d = roslices.Clone(data.RO)
	if !d.inOrder() {
		t.Errorf("Stable shuffled sorted %d ints (stability)", n)
	}

	// sorted reversed
	for i := 0; i < len(d); i++ {
		d[i].a = len(d) - i
	}
	d.initB()
	SortStableFunc(&data, intPairLess)
	if !roslices.IsSortedFunc(data.RO, intPairLess) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	d = roslices.Clone(data.RO)
	if !d.inOrder() {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}
}

var insertTests = []struct {
	s    Slice[int]
	i    int
	add  []int
	want Slice[int]
}{
	{
		CopyOnWrite([]int{1, 2, 3}),
		0,
		[]int{4},
		CopyOnWrite([]int{4, 1, 2, 3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		1,
		[]int{4},
		CopyOnWrite([]int{1, 4, 2, 3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		3,
		[]int{4},
		CopyOnWrite([]int{1, 2, 3, 4}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		2,
		[]int{4, 5},
		CopyOnWrite([]int{1, 2, 4, 5, 3}),
	},
}

func TestInsert(t *testing.T) {
	s := CopyOnWrite([]int{1, 2, 3})
	if got := Insert(s, 0); !roslices.Equal(got.RO, s.RO) {
		t.Errorf("Insert(%v, 0) = %v, want %v", s, got, s)
	}
	for _, test := range insertTests {
		if got := Insert(test.s, test.i, test.add...); !roslices.Equal(got.RO, test.want.RO) {
			t.Errorf("Insert(%v, %d, %v...) = %v, want %v", test.s, test.i, test.add, got, test.want)
		}
	}
}

var deleteTests = []struct {
	s    Slice[int]
	i, j int
	want Slice[int]
}{
	{
		CopyOnWrite([]int{1, 2, 3}),
		0,
		0,
		CopyOnWrite([]int{1, 2, 3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		0,
		1,
		CopyOnWrite([]int{2, 3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		3,
		3,
		CopyOnWrite([]int{1, 2, 3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		0,
		2,
		CopyOnWrite([]int{3}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		0,
		3,
		CopyOnWrite([]int{}),
	},
}

func TestDelete(t *testing.T) {
	for _, test := range deleteTests {
		if got := Delete(test.s, test.i, test.j); !roslices.Equal(got.RO, test.want.RO) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}

var compactTests = []struct {
	s    Slice[int]
	want Slice[int]
}{
	{
		CopyOnWrite[int](nil),
		CopyOnWrite[int](nil),
	},
	{
		CopyOnWrite([]int{1}),
		CopyOnWrite([]int{1}),
	},
	{
		CopyOnWrite([]int{1, 2, 3}),
		CopyOnWrite([]int{1, 2, 3}),
	},
	{
		CopyOnWrite([]int{1, 1, 2}),
		CopyOnWrite([]int{1, 2}),
	},
	{
		CopyOnWrite([]int{1, 2, 1}),
		CopyOnWrite([]int{1, 2, 1}),
	},
	{
		CopyOnWrite([]int{1, 2, 2, 3, 3, 4}),
		CopyOnWrite([]int{1, 2, 3, 4}),
	},
}

func TestCompact(t *testing.T) {
	for _, test := range compactTests {
		if got := Compact(test.s); !roslices.Equal(got.RO, test.want.RO) {
			t.Errorf("Compact(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

func TestCompactFunc(t *testing.T) {
	for _, test := range compactTests {
		if got := CompactFunc(test.s, equal[int]); !roslices.Equal(got.RO, test.want.RO) {
			t.Errorf("CompactFunc(%v, equal[int]) = %v, want %v", test.s, got, test.want)
		}
	}

	s1 := CopyOnWrite([]string{"a", "a", "A", "B", "b"})
	want := CopyOnWrite([]string{"a", "B"})
	if got := CompactFunc(s1, strings.EqualFold); !roslices.Equal(got.RO, want.RO) {
		t.Errorf("CompactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}

func TestGrow(t *testing.T) {
	s1 := CopyOnWrite([]int{1, 2, 3})
	s2 := Grow(s1, 1000)
	if !roslices.Equal(s1.RO, s2.RO) {
		t.Errorf("Grow(%v) = %v, want %v", s1, s2, s1)
	}
	if s2.RO.Cap() < 1000+s1.RO.Len() {
		t.Errorf("after Grow(%v) cap = %d, want >= %d", s1, s2.RO.Cap(), 1000+s1.RO.Len())
	}
}

func TestClip(t *testing.T) {
	orig := []int{1, 2, 3, 4, 5, 6}[:3]
	s1 := CopyOnWrite(orig)
	if s1.RO.Len() != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, s1.RO.Len())
	}
	if s1.RO.Cap() < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", orig, s1.RO.Cap())
	}
	s2 := Clip(s1)
	if !roslices.Equal(s1.RO, s2.RO) {
		t.Errorf("Clip(%v) = %v, want %v", s1, s2, s1)
	}
	if s2.RO.Cap() != 3 {
		t.Errorf("cap(Clip(%v)) = %d, want 3", orig, s2.RO.Cap())
	}
}

// These benchmarks compare sorting a large slice of int with sort.Ints vs.
// Sort
func makeRandomInts(n int) []int {
	rand.Seed(42)
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = rand.Intn(n)
	}
	return ints
}

const N = 100_000

func BenchmarkSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		sort.Ints(ints)
	}
}

func BenchmarkSlicesSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := CopyOnWrite(makeRandomInts(N))
		b.StartTimer()
		Sort(&ints)
	}
}

// Since we're benchmarking these sorts against each other, make sure that they
// generate similar results.
func TestIntSorts(t *testing.T) {
	ints := makeRandomInts(200)
	ints2 := CopyOnWrite(ints)

	sort.Ints(ints)
	Sort(&ints2)

	for i := range ints {
		if ints[i] != ints2.RO.Index(i) {
			t.Fatalf("ints2 mismatch at %d; %d != %d", i, ints[i], ints2.RO.Index(i))
		}
	}
}

// The following is a benchmark for sorting strings.

// makeRandomStrings generates n random strings with alphabetic runes of
// varying lengths.
func makeRandomStrings(n int) []string {
	rand.Seed(42)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		slen := 2 + rand.Intn(50)
		for j := 0; j < slen; j++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		ss[i] = sb.String()
	}
	return ss
}

func TestStringSorts(t *testing.T) {
	ss := makeRandomStrings(200)
	ss2 := CopyOnWrite(ss)

	sort.Strings(ss)
	Sort(&ss2)

	for i := range ss {
		if ss[i] != ss2.RO.Index(i) {
			t.Fatalf("ss2 mismatch at %d; %s != %s", i, ss[i], ss2.RO.Index(i))
		}
	}
}

func BenchmarkSortStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		sort.Strings(ss)
	}
}

func BenchmarkSlicesSortStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := CopyOnWrite(makeRandomStrings(N))
		b.StartTimer()
		Sort(&ss)
	}
}

// These benchmarks compare sorting a slice of structs with sort.Sort vs.
// slices.SortFunc.
type myStruct struct {
	a, b, c, d string
	n          int
}

type myStructs []*myStruct

func (s myStructs) Len() int           { return len(s) }
func (s myStructs) Less(i, j int) bool { return s[i].n < s[j].n }
func (s myStructs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func makeRandomStructs(n int) myStructs {
	rand.Seed(42)
	structs := make([]*myStruct, n)
	for i := 0; i < n; i++ {
		structs[i] = &myStruct{n: rand.Intn(n)}
	}
	return structs
}

func TestStructSorts(t *testing.T) {
	ss := makeRandomStructs(200)
	ss2orig := make([]*myStruct, len(ss))
	for i := range ss {
		ss2orig[i] = &myStruct{n: ss[i].n}
	}
	ss2 := CopyOnWrite(ss2orig)

	sort.Sort(ss)
	SortFunc(&ss2, func(a, b *myStruct) bool { return a.n < b.n })

	for i := range ss {
		if *ss[i] != *ss2.RO.Index(i) {
			t.Fatalf("ints2 mismatch at %d; %v != %v", i, *ss[i], *ss2.RO.Index(i))
		}
	}
}

func BenchmarkSortStructs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStructs(N)
		b.StartTimer()
		sort.Sort(ss)
	}
}

func BenchmarkSortFuncStructs(b *testing.B) {
	lessFunc := func(a, b *myStruct) bool { return a.n < b.n }
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := CopyOnWrite(makeRandomStructs(N))
		b.StartTimer()
		SortFunc(&ss, lessFunc)
	}
}
