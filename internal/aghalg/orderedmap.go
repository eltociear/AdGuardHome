package aghalg

import (
	"golang.org/x/exp/slices"
)

// OrderedMap is the implementation of the ordered map data structure.
type OrderedMap[K comparable, T any] struct {
	vals map[K]T
	cmp  func(a, b K) int
	keys []K
}

// NewOrderedMap initializes the new instance of ordered map.  cmp is a sort
// function.
func NewOrderedMap[K comparable, T any](cmp func(a, b K) int) OrderedMap[K, T] {
	return OrderedMap[K, T]{
		vals: make(map[K]T),
		cmp:  cmp,
	}
}

// Add adds val with key to the ordered map.
func (m *OrderedMap[K, T]) Add(key K, val T) {
	i, has := slices.BinarySearchFunc(m.keys, key, m.cmp)
	if has {
		m.keys[i] = key
		m.vals[key] = val

		return
	}

	m.keys = slices.Insert(m.keys, i, key)
	m.vals[key] = val
}

// Del removes the value by key from the ordered map.
func (m *OrderedMap[K, T]) Del(key K) {
	i, has := slices.BinarySearchFunc(m.keys, key, m.cmp)
	if !has {
		return
	}

	m.keys = slices.Delete(m.keys, i, 1)
	delete(m.vals, key)
}

// Range calls cb for each element of the map.  If cb returns false it stops.
func (m *OrderedMap[K, T]) Range(cb func(K, T) bool) {
	for _, k := range m.keys {
		if !cb(k, m.vals[k]) {
			return
		}
	}
}
