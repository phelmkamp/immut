package cowmaps

import (
	"github.com/phelmkamp/immut/romaps"
	"golang.org/x/exp/maps"
)

// Doer defines a method for doing an operation on a map.
type Doer[K comparable, V any] interface {
	do(map[K]V) map[K]V
}

// doerFunc is a func that implements Doer.
type doerFunc[K comparable, V any] func(map[K]V) map[K]V

func (f doerFunc[K, V]) do(m map[K]V) map[K]V {
	return f(m)
}

// DoAll does all the supplied operations on the map with minimal reallocation.
// The initial capacity of the reallocated map is cap (or len(m) if cap is not sufficient).
// Note: The underlying map is cloned before the write-operations are performed.
func DoAll[K comparable, V any](m *Map[K, V], cap int, ops ...Doer[K, V]) {
	ro := m.RO
	if cap < ro.Len() {
		cap = ro.Len()
	}
	m2 := clone(ro, cap)
	for _, op := range ops {
		m2 = op.do(m2)
	}
	m.RO = romaps.Freeze(m2)
}

// DoCopy returns the maps.Copy operation.
func DoCopy[K comparable, V any](src map[K]V) Doer[K, V] {
	return doerFunc[K, V](func(m map[K]V) map[K]V {
		maps.Copy(m, src)
		return m
	})
}

// DoDelete returns the delete operation.
func DoDelete[K comparable, V any](k K) Doer[K, V] {
	return doerFunc[K, V](func(m map[K]V) map[K]V {
		delete(m, k)
		return m
	})
}

// DoDeleteFunc returns the maps.DeleteFunc operation.
func DoDeleteFunc[K comparable, V any](del func(K, V) bool) Doer[K, V] {
	return doerFunc[K, V](func(m map[K]V) map[K]V {
		maps.DeleteFunc(m, del)
		return m
	})
}

// DoSetIndex returns the set index operation.
func DoSetIndex[K comparable, V any](k K, v V) Doer[K, V] {
	return doerFunc[K, V](func(m map[K]V) map[K]V {
		m[k] = v
		return m
	})
}
