// Package hashset implements a set data structure based on hash table.
package hashset

import "iter"

// HashSet holds a set of unique elements
type HashSet[K comparable] struct {
	data map[K]struct{}
}

// Equal reports whether two sets contain the same elements
func Equal[K comparable](s1, s2 *HashSet[K]) bool {
	if len(s1.data) != len(s2.data) {
		return false
	}
	for v := range s1.data {
		if _, ok := s2.data[v]; !ok {
			return false
		}
	}
	return true
}

// New creates and initialize a new [HashSet]
func New[K comparable]() *HashSet[K] {
	s := &HashSet[K]{
		data: make(map[K]struct{}),
	}
	return s
}

// Contains report whether an element exists in the set
func (s *HashSet[K]) Contains(key K) bool {
	_, ok := s.data[key]
	return ok
}

// Empty reports whether the set is empty
func (s *HashSet[K]) Empty() bool {
	return len(s.data) == 0
}

// Len returns the number of elements in the set
func (s *HashSet[K]) Len() int {
	return len(s.data)
}

// Clear removes all elements from the set
func (s *HashSet[K]) Clear() {
	clear(s.data)
}

// Add adds an element to the set
// If an element already exists in the set, it is ignored.
func (s *HashSet[K]) Add(key K) {
	s.data[key] = struct{}{}
}

// Remove removes an element from the set.
// If an element does not exist in the set, it is ignored.
func (s *HashSet[K]) Remove(key K) {
	delete(s.data, key)
}

// All is an iterator over the elements in the set
func (s *HashSet[K]) All() iter.Seq[K] {
	return func(yield func(K) bool) {
		for v := range s.data {
			if !yield(v) {
				return
			}
		}
	}
}

// Collect creates a new set from an iterator
func Collect[K comparable](seq iter.Seq[K]) *HashSet[K] {
	s := New[K]()
	for v := range seq {
		s.data[v] = struct{}{}
	}
	return s
}

// Union returns a new set that contains all elements from two sets
func Union[K comparable](s1, s2 *HashSet[K]) *HashSet[K] {
	result := New[K]()
	for v := range s1.data {
		result.data[v] = struct{}{}
	}
	for v := range s2.data {
		result.data[v] = struct{}{}
	}
	return result
}

// Intersection returns a new set that contains common elements from two sets
func Intersection[K comparable](s1, s2 *HashSet[K]) *HashSet[K] {
	result := New[K]()
	for v := range s1.data {
		if _, ok := s2.data[v]; ok {
			result.data[v] = struct{}{}
		}
	}

	return result
}

// Difference returns a new set that contains elements
// that are in the first set but not in the second set
func Difference[K comparable](s1, s2 *HashSet[K]) *HashSet[K] {
	result := New[K]()
	for v := range s1.data {
		if _, ok := s2.data[v]; !ok {
			result.data[v] = struct{}{}
		}
	}

	return result
}
