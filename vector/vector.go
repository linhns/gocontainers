// Package vector provides a simple vector implementation
// based on Go slices.
package vector

import (
	"iter"
	"slices"
)

// Vector represent a growable collection of elements
// that are stored contiguously in memory.
type Vector[T any] struct {
	data []T
}

// New creates and initializes a new [Vector]
func New[T any]() *Vector[T] {
	return &Vector[T]{
		data: make([]T, 0),
	}
}

// NewWithCapacity creates and initializes a new [Vector]
// with a specified capacity
func NewWithCapacity[T any](capacity int) *Vector[T] {
	return &Vector[T]{
		data: make([]T, 0, capacity),
	}
}

// Push adds an element to the end of the vector
func (v *Vector[T]) Push(val T) {
	v.data = append(v.data, val)
}

// Pop removes and returns the last element from the vector.
// If the vector is empty, it returns a zero value and false.
func (v *Vector[T]) Pop() (T, bool) {
	var zero T
	if len(v.data) == 0 {
		return zero, false
	}
	idx := len(v.data) - 1
	val := v.data[idx]
	v.data = v.data[:idx]
	return val, true
}

// Len returns the number of elements in the vector
func (v *Vector[T]) Len() int {
	return len(v.data)
}

// Capacity returns the number of elements that the vector
// can hold without further allocation
func (v *Vector[T]) Capacity() int {
	return cap(v.data)
}

// Empty reports whether the vector is empty
func (v *Vector[T]) Empty() bool {
	return len(v.data) == 0
}

// Get returns the element at the specified index.
//
// If the index is out of range, it returns a zero value and false.
// Otherwise, it returns the element and true.
func (v *Vector[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(v.data) {
		return zero, false
	}

	return v.data[index], true
}

// Set sets the element at the specified index.
// If the index is out of range, it returns false.
func (v *Vector[T]) Set(index int, value T) bool {
	if index < 0 || index >= len(v.data) {
		return false
	}

	v.data[index] = value
	return true
}

// Front returns the first element of the vector.
//
// If the vector is empty, it returns a zero value and false.
// Otherwise, it returns the first element and true.
func (v *Vector[T]) Front() (T, bool) {
	return v.Get(0)
}

// Back returns the last element of the vector.
//
// If the vector is empty, it returns a zero value and false.
// Otherwise, it returns the last element and true.
func (v *Vector[T]) Back() (T, bool) {
	return v.Get(v.Len() - 1)
}

// Clear removes all elements from the vector
func (v *Vector[T]) Clear() {
	v.data = v.data[:0]
}

// Reserve increases the capacity of the vector to
// the specified value.
//
// If the specified capacity is less than or equal to the current capacity,
// this is a no-op.
func (v *Vector[T]) Reserve(capacity int) {
	if capacity <= cap(v.data) {
		return
	}
	data := make([]T, len(v.data), capacity)
	copy(data, v.data)
	v.data = data
}

// ShrinkToFit reduces the capacity of the vector as much as possible.
func (v *Vector[T]) ShrinkToFit() {
	if len(v.data) == cap(v.data) {
		return
	}

	sz := len(v.data)
	data := make([]T, sz)
	copy(data, v.data)
	v.data = data
}

// Resize changes the size of the vector to the specified size.
//
// If the specified size is less than the current size, the vector
// is reduced to the first size elements.
//
// If the specified size is greater than the current size, the vector
// additional zero-valued elements are appended.
func (v *Vector[T]) Resize(size int) {
	if size < len(v.data) {
		v.data = v.data[:size]
		return
	}
	var zero T
	for i := len(v.data); i < size; i++ {
		v.data = append(v.data, zero)
	}
}

func (v *Vector[T]) Insert(index int, value T) {
	if index < 0 || index > len(v.data) {
		return
	}

	v.data = append(v.data, value)
	copy(v.data[index+1:], v.data[index:])
	v.data[index] = value
}

// Equal reports whether two vectors are equal.
func Equal[T comparable](v1, v2 *Vector[T]) bool {
	if v1.Len() != v2.Len() {
		return false
	}

	for i := 0; i < v1.Len(); i++ {
		if v1.data[i] != v2.data[i] {
			return false
		}
	}

	return true
}

// Values returns an iterator that yields the vector elements in order.
func (v *Vector[T]) Values() iter.Seq[T] {
	return slices.Values(v.data)
}

// All returns an iterator over index-value pairs in the vector.
func (v *Vector[T]) All() iter.Seq2[int, T] {
	return slices.All(v.data)
}

// Collect collects values from an iterator and returns a new vector.
func Collect[T any](seq iter.Seq[T]) *Vector[T] {
	v := New[T]()
	for val := range seq {
		v.Push(val)
	}
	return v
}
