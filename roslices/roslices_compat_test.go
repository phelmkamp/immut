package roslices

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

// copied from https://cs.opensource.google/go/x/exp/+/master:slices/
// to ensure compatibility

var equalIntTests = []struct {
	s1, s2 Slice[int]
	want   bool
}{
	{
		Freeze([]int{1}),
		Freeze[int](nil),
		false,
	},
	{
		Freeze([]int{}),
		Freeze[int](nil),
		true,
	},
	{
		Freeze([]int{1, 2, 3}),
		Freeze([]int{1, 2, 3}),
		true,
	},
	{
		Freeze([]int{1, 2, 3}),
		Freeze([]int{1, 2, 3, 4}),
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       Slice[float64]
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		Freeze([]float64{1, 2}),
		Freeze([]float64{1, 2}),
		true,
		true,
	},
	{
		Freeze([]float64{1, 2, math.NaN()}),
		Freeze([]float64{1, 2, math.NaN()}),
		false,
		true,
	},
}

func TestEqual(t *testing.T) {
	for _, test := range equalIntTests {
		if got := Equal(test.s1, test.s2); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := Equal(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
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

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne[E constraints.Integer](v1, v2 E) bool {
	return v1 == v2+1 || v1 == v2-1
}

func TestEqualFunc(t *testing.T) {
	for _, test := range equalIntTests {
		if got := EqualFunc(test.s1, test.s2, equal[int]); got != test.want {
			t.Errorf("EqualFunc(%v, %v, equal[int]) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := EqualFunc(test.s1, test.s2, equal[float64]); got != test.wantEqual {
			t.Errorf("Equal(%v, %v, equal[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
		if got := EqualFunc(test.s1, test.s2, equalNaN[float64]); got != test.wantEqualNaN {
			t.Errorf("Equal(%v, %v, equalNaN[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqualNaN)
		}
	}

	s1 := Freeze([]int{1, 2, 3})
	s2 := Freeze([]int{2, 3, 4})
	if EqualFunc(s1, s1, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}
	if !EqualFunc(s1, s2, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}

	s3 := Freeze([]string{"a", "b", "c"})
	s4 := Freeze([]string{"A", "B", "C"})
	if !EqualFunc(s3, s4, strings.EqualFold) {
		t.Errorf("EqualFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}

	cmpIntString := func(v1 int, v2 string) bool {
		return string(rune(v1)-1+'a') == v2
	}
	if !EqualFunc(s1, s3, cmpIntString) {
		t.Errorf("EqualFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
}

var compareIntTests = []struct {
	s1, s2 Slice[int]
	want   int
}{
	{
		Freeze([]int{1, 2, 3}),
		Freeze([]int{1, 2, 3, 4}),
		-1,
	},
	{
		Freeze([]int{1, 2, 3, 4}),
		Freeze([]int{1, 2, 3}),
		+1,
	},
	{
		Freeze([]int{1, 2, 3}),
		Freeze([]int{1, 4, 3}),
		-1,
	},
	{
		Freeze([]int{1, 4, 3}),
		Freeze([]int{1, 2, 3}),
		+1,
	},
}

var compareFloatTests = []struct {
	s1, s2 Slice[float64]
	want   int
}{
	{
		Freeze([]float64{1, 2, math.NaN()}),
		Freeze([]float64{1, 2, math.NaN()}),
		0,
	},
	{
		Freeze([]float64{1, math.NaN(), 3}),
		Freeze([]float64{1, math.NaN(), 4}),
		-1,
	},
	{
		Freeze([]float64{1, math.NaN(), 3}),
		Freeze([]float64{1, 2, 4}),
		-1,
	},
	{
		Freeze([]float64{1, math.NaN(), 3}),
		Freeze([]float64{1, 2, math.NaN()}),
		0,
	},
	{
		Freeze([]float64{1, math.NaN(), 3, 4}),
		Freeze([]float64{1, 2, math.NaN()}),
		+1,
	},
}

func TestCompare(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.wantEqualNaN {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqualNaN))
		}
	}

	for _, test := range compareIntTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
}

func equalToCmp[T comparable](eq func(T, T) bool) func(T, T) int {
	return func(v1, v2 T) int {
		if eq(v1, v2) {
			return 0
		}
		return 1
	}
}

func cmp[T constraints.Ordered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	} else {
		return 0
	}
}

func TestCompareFunc(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := CompareFunc(test.s1, test.s2, equalToCmp(equal[int])); (got == 0) != test.want {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[int])) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := CompareFunc(test.s1, test.s2, equalToCmp(equal[float64])); (got == 0) != test.wantEqual {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[float64])) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqual))
		}
	}

	for _, test := range compareIntTests {
		if got := CompareFunc(test.s1, test.s2, cmp[int]); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[int]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := CompareFunc(test.s1, test.s2, cmp[float64]); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[float64]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	s1 := Freeze([]int{1, 2, 3})
	s2 := Freeze([]int{2, 3, 4})
	if got := CompareFunc(s1, s2, equalToCmp(offByOne[int])); got != 0 {
		t.Errorf("CompareFunc(%v, %v, offByOne) = %d, want 0", s1, s2, got)
	}

	s3 := Freeze([]string{"a", "b", "c"})
	s4 := Freeze([]string{"A", "B", "C"})
	if got := CompareFunc(s3, s4, strings.Compare); got != 1 {
		t.Errorf("CompareFunc(%v, %v, strings.Compare) = %d, want 1", s3, s4, got)
	}

	compareLower := func(v1, v2 string) int {
		return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
	}
	if got := CompareFunc(s3, s4, compareLower); got != 0 {
		t.Errorf("CompareFunc(%v, %v, compareLower) = %d, want 0", s3, s4, got)
	}

	cmpIntString := func(v1 int, v2 string) int {
		return strings.Compare(string(rune(v1)-1+'a'), v2)
	}
	if got := CompareFunc(s1, s3, cmpIntString); got != 0 {
		t.Errorf("CompareFunc(%v, %v, cmpIntString) = %d, want 0", s1, s3, got)
	}
}

var indexTests = []struct {
	s    Slice[int]
	v    int
	want int
}{
	{
		Freeze[int](nil),
		0,
		-1,
	},
	{
		Freeze([]int{}),
		0,
		-1,
	},
	{
		Freeze([]int{1, 2, 3}),
		2,
		1,
	},
	{
		Freeze([]int{1, 2, 2, 3}),
		2,
		1,
	},
	{
		Freeze([]int{1, 2, 3, 2}),
		2,
		1,
	},
}

func TestIndex(t *testing.T) {
	for _, test := range indexTests {
		if got := Index(test.s, test.v); got != test.want {
			t.Errorf("Index(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}

func equalToIndex[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}

func TestIndexFunc(t *testing.T) {
	for _, test := range indexTests {
		if got := IndexFunc(test.s, equalToIndex(equal[int], test.v)); got != test.want {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", test.s, test.v, got, test.want)
		}
	}

	s1 := Freeze([]string{"hi", "HI"})
	if got := IndexFunc(s1, equalToIndex(equal[string], "HI")); got != 1 {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, "HI", got, 1)
	}
	if got := IndexFunc(s1, equalToIndex(strings.EqualFold, "HI")); got != 0 {
		t.Errorf("IndexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, "HI", got, 0)
	}
}

func TestContains(t *testing.T) {
	for _, test := range indexTests {
		if got := Contains(test.s, test.v); got != (test.want != -1) {
			t.Errorf("Contains(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}

func TestClone(t *testing.T) {
	s0 := []int{1, 2, 3}
	s1 := Freeze(s0)
	s2 := Freeze(Clone(s1))
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %v", s1, s2, s1)
	}
	s0[0] = 4
	want := Freeze([]int{1, 2, 3})
	if !Equal(s2, want) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}
	if got := Clone(Freeze([]int(nil))); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone(Freeze(s0[:0])); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s0[:0], got, s0[:0])
	}
}

func TestBinarySearch(t *testing.T) {
	str1 := Freeze([]string{"foo"})
	str2 := Freeze([]string{"ab", "ca"})
	str3 := Freeze([]string{"mo", "qo", "vo"})
	str4 := Freeze([]string{"ab", "ad", "ca", "xy"})

	// slice with repeating elements
	strRepeats := Freeze([]string{"ba", "ca", "da", "da", "da", "ka", "ma", "ma", "ta"})

	// slice with all element equal
	strSame := Freeze([]string{"xx", "xx", "xx"})

	tests := []struct {
		data      Slice[string]
		target    string
		wantPos   int
		wantFound bool
	}{
		{Freeze([]string{}), "foo", 0, false},
		{Freeze([]string{}), "", 0, false},

		{str1, "foo", 0, true},
		{str1, "bar", 0, false},
		{str1, "zx", 1, false},

		{str2, "aa", 0, false},
		{str2, "ab", 0, true},
		{str2, "ad", 1, false},
		{str2, "ca", 1, true},
		{str2, "ra", 2, false},

		{str3, "bb", 0, false},
		{str3, "mo", 0, true},
		{str3, "nb", 1, false},
		{str3, "qo", 1, true},
		{str3, "tr", 2, false},
		{str3, "vo", 2, true},
		{str3, "xr", 3, false},

		{str4, "aa", 0, false},
		{str4, "ab", 0, true},
		{str4, "ac", 1, false},
		{str4, "ad", 1, true},
		{str4, "ax", 2, false},
		{str4, "ca", 2, true},
		{str4, "cc", 3, false},
		{str4, "dd", 3, false},
		{str4, "xy", 3, true},
		{str4, "zz", 4, false},

		{strRepeats, "da", 2, true},
		{strRepeats, "db", 5, false},
		{strRepeats, "ma", 6, true},
		{strRepeats, "mb", 8, false},

		{strSame, "xx", 0, true},
		{strSame, "ab", 0, false},
		{strSame, "zz", 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			{
				pos, found := BinarySearch(tt.data, tt.target)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}

			{
				pos, found := BinarySearchFunc(tt.data, tt.target, strings.Compare)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}
		})
	}
}

func TestBinarySearchInts(t *testing.T) {
	data := Freeze([]int{20, 30, 40, 50, 60, 70, 80, 90})
	tests := []struct {
		target    int
		wantPos   int
		wantFound bool
	}{
		{20, 0, true},
		{23, 1, false},
		{43, 3, false},
		{80, 6, true},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.target), func(t *testing.T) {
			{
				pos, found := BinarySearch(data, tt.target)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}

			{
				cmp := func(a, b int) int {
					return a - b
				}
				pos, found := BinarySearchFunc(data, tt.target, cmp)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("BinarySearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
			}
		})
	}
}
