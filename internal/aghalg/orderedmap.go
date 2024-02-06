package aghalg

import (
	"golang.org/x/exp/slices"
)

// SortedMap is a map that keeps elements in order with internal sorting
// function.
type SortedMap[K comparable, V any] struct {
	vals map[K]V
	cmp  func(a, b K) (res int)
	keys []K
}

// NewSortedMap initializes the new instance of sorted map.  cmp is a sort
// function to keep elements in order.
//
// TODO(s.chzhen):  Use cmp.Compare in Go 1.21.
func NewSortedMap[K comparable, V any](cmp func(a, b K) (res int)) SortedMap[K, V] {
	return SortedMap[K, V]{
		vals: make(map[K]V),
		cmp:  cmp,
	}
}

// Set adds val with key to the sorted map.
func (m *SortedMap[K, V]) Set(key K, val V) {
	i, has := slices.BinarySearchFunc(m.keys, key, m.cmp)
	if has {
		m.keys[i] = key
		m.vals[key] = val

		return
	}

	m.keys = slices.Insert(m.keys, i, key)
	m.vals[key] = val
}

// Get returns val by key from the sorted map.
func (m *SortedMap[K, V]) Get(key K) (val V) {
	return m.vals[key]
}

// Del removes the value by key from the sorted map.
func (m *SortedMap[K, V]) Del(key K) {
	i, has := slices.BinarySearchFunc(m.keys, key, m.cmp)
	if has {
		m.keys = slices.Delete(m.keys, i, i+1)
		delete(m.vals, key)
	}
}

// Clear removes all elements from the sorted map.
func (m *SortedMap[K, V]) Clear() {
	// TODO(s.chzhen):  Use built-in clear in Go 1.21.
	m.keys = nil
	m.vals = make(map[K]V)
}

// Range calls cb for each element of the map, sorted by m.cmp.  If cb returns
// false it stops.
func (m *SortedMap[K, V]) Range(cb func(K, V) (cont bool)) {
	for _, k := range m.keys {
		if !cb(k, m.vals[k]) {
			return
		}
	}
}
