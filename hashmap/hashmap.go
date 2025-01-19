// Package hashmap wraps the Go built-in map and provides a more intuitive API.
package hashmap

import (
	"iter"
	"maps"
)

// HashMap is a generic hash table (map).
type HashMap[K comparable, V any] struct {
	data map[K]V
}

// New creates and initialize a new [HashMap].
func New[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		data: make(map[K]V),
	}
}

// Insert inserts a key-value pair into the map.
//
// If the map does not contain the key, it will be added.
//
// If the map already contains the key, the value will be updated.
func (m *HashMap[K, V]) Insert(key K, value V) {
	m.data[key] = value
}

// Get retrieves the value associated with the key. If the key does not exist,
// it returns the zero value of the value type and false.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
	value, ok := m.data[key]
	if !ok {
		var zero V
		return zero, false
	}
	return value, true
}

// Contains reports whether the map contains the key.
func (m *HashMap[K, V]) Contains(key K) bool {
	_, ok := m.data[key]
	return ok
}

// Remove removes the key-value pair from the map. If the key does not exist,
// this is a no-op.
func (m *HashMap[K, V]) Remove(key K) {
	delete(m.data, key)
}

// Clear removes all key-value pairs from the map.
func (m *HashMap[K, V]) Clear() {
	clear(m.data)
}

// Len returns the number of key-value pairs in the map.
func (m *HashMap[K, V]) Len() int {
	return len(m.data)
}

// Empty reports whether the map is empty.
func (m *HashMap[K, V]) Empty() bool {
	return len(m.data) == 0
}

// Keys returns an iterator over keys in the map.
// The iteration order is unspecified and not guaranteed
// to remain the same between calls.
func (m *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.data)
}

// Values returns an iterator over values in the map.
// The iteration order is unspecified and not guaranteed
// to remain the same between calls.
func (m *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.data)
}

// All returns an iterator over key-value pairs in the map.
// The iteration order is unspecified and not guaranteed
// to remain the same between calls.
func (m *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.data)
}

// Collect collects key-value pairs from an iterator and returns a new map.
func Collect[K comparable, V any](seq iter.Seq2[K, V]) *HashMap[K, V] {
	m := New[K, V]()
	for k, v := range seq {
		m.data[k] = v
	}
	return m
}
