package aghalg

import (
	"golang.org/x/exp/slices"
)

// OrderedMap is the implementation of the ordered map data structure.
type OrderedMap[K comparable, V any] struct {
	vals map[K]V
	cmp  func(a, b K) int
	keys []K
}

// NewOrderedMap initializes the new instance of ordered map.  cmp is a sort
// function.
//
// TODO(s.chzhen):  Use cmp.Compare in Go 1.21
func NewOrderedMap[K comparable, V any](cmp func(a, b K) int) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		vals: make(map[K]V),
		cmp:  cmp,
	}
}

// Set adds val with key to the ordered map.
func (m *OrderedMap[K, V]) Set(key K, val V) {
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
func (m *OrderedMap[K, V]) Del(key K) {
	i, has := slices.BinarySearchFunc(m.keys, key, m.cmp)
	if has {
		m.keys = slices.Delete(m.keys, i, 1)
		delete(m.vals, key)
	}
}

// Range calls cb for each element of the map.  If cb returns false it stops.
func (m *OrderedMap[K, V]) Range(cb func(K, V) (cont bool)) {
	for _, k := range m.keys {
		if !cb(k, m.vals[k]) {
			return
		}
	}
}
