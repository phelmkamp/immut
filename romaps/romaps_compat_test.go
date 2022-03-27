package romaps

import (
	"math"
	"sort"
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

// copied from https://cs.opensource.google/go/x/exp/+/master:slices/
// to ensure compatibility

var m1 = Freeze(map[int]int{1: 2, 2: 4, 4: 8, 8: 16})
var m2 = Freeze(map[int]string{1: "2", 2: "4", 4: "8", 8: "16"})

func TestKeys(t *testing.T) {
	want := []int{1, 2, 4, 8}

	got1 := Keys(m1)
	sort.Ints(got1)
	if !slices.Equal(got1, want) {
		t.Errorf("Keys(%v) = %v, want %v", m1, got1, want)
	}

	got2 := Keys(m2)
	sort.Ints(got2)
	if !slices.Equal(got2, want) {
		t.Errorf("Keys(%v) = %v, want %v", m2, got2, want)
	}
}

func TestValues(t *testing.T) {
	got1 := Values(m1)
	want1 := []int{2, 4, 8, 16}
	sort.Ints(got1)
	if !slices.Equal(got1, want1) {
		t.Errorf("Values(%v) = %v, want %v", m1, got1, want1)
	}

	got2 := Values(m2)
	want2 := []string{"16", "2", "4", "8"}
	sort.Strings(got2)
	if !slices.Equal(got2, want2) {
		t.Errorf("Values(%v) = %v, want %v", m2, got2, want2)
	}
}

func TestEqual(t *testing.T) {
	if !Equal(m1, m1) {
		t.Errorf("Equal(%v, %v) = false, want true", m1, m1)
	}
	if Equal(m1, Freeze((map[int]int)(nil))) {
		t.Errorf("Equal(%v, nil) = true, want false", m1)
	}
	if Equal(Freeze((map[int]int)(nil)), m1) {
		t.Errorf("Equal(nil, %v) = true, want false", m1)
	}
	if !Equal[int, int](Freeze((map[int]int)(nil)), Freeze((map[int]int)(nil))) {
		t.Error("Equal(nil, nil) = false, want true")
	}
	if ms := Freeze(map[int]int{1: 2}); Equal(m1, ms) {
		t.Errorf("Equal(%v, %v) = true, want false", m1, ms)
	}

	// Comparing NaN for equality is expected to fail.
	mf := Freeze(map[int]float64{1: 0, 2: math.NaN()})
	if Equal(mf, mf) {
		t.Errorf("Equal(%v, %v) = true, want false", mf, mf)
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// equalStr compares ints and strings.
func equalIntStr(v1 int, v2 string) bool {
	return strconv.Itoa(v1) == v2
}

func TestEqualFunc(t *testing.T) {
	if !EqualFunc(m1, m1, equal[int]) {
		t.Errorf("EqualFunc(%v, %v, equal) = false, want true", m1, m1)
	}
	if EqualFunc(m1, Freeze((map[int]int)(nil)), equal[int]) {
		t.Errorf("EqualFunc(%v, nil, equal) = true, want false", m1)
	}
	if EqualFunc(Freeze((map[int]int)(nil)), m1, equal[int]) {
		t.Errorf("EqualFunc(nil, %v, equal) = true, want false", m1)
	}
	if !EqualFunc[int, int](Freeze((map[int]int)(nil)), Freeze((map[int]int)(nil)), equal[int]) {
		t.Error("EqualFunc(nil, nil, equal) = false, want true")
	}
	if ms := Freeze(map[int]int{1: 2}); EqualFunc(m1, ms, equal[int]) {
		t.Errorf("EqualFunc(%v, %v, equal) = true, want false", m1, ms)
	}

	// Comparing NaN for equality is expected to fail.
	mf := Freeze(map[int]float64{1: 0, 2: math.NaN()})
	if EqualFunc(mf, mf, equal[float64]) {
		t.Errorf("EqualFunc(%v, %v, equal) = true, want false", mf, mf)
	}
	// But it should succeed using equalNaN.
	if !EqualFunc(mf, mf, equalNaN[float64]) {
		t.Errorf("EqualFunc(%v, %v, equalNaN) = false, want true", mf, mf)
	}

	if !EqualFunc(m1, m2, equalIntStr) {
		t.Errorf("EqualFunc(%v, %v, equalIntStr) = false, want true", m1, m2)
	}
}

func TestClone(t *testing.T) {
	mc := Clone(m1)
	if !Equal(Freeze(mc), m1) {
		t.Errorf("Clone(%v) = %v, want %v", m1, mc, m1)
	}
	mc[16] = 32
	if Equal(Freeze(mc), m1) {
		t.Errorf("Equal(%v, %v) = true, want false", mc, m1)
	}
}
