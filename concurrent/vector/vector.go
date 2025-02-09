// Package vector provides a simple vector implementation.
package vector

import (
	"iter"
	"slices"
	"sync"
)

// Vector represent a growable collection of elements
// that are stored contiguously in memory.
type Vector[T any] struct {
	data []T
	mu   sync.RWMutex
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

// Of constructs a new [Vector] with initial values
func Of[T any](v ...T) *Vector[T] {
	return &Vector[T]{
		data: v,
	}
}

// PushBack adds an element to the end of the vector
func (v *Vector[T]) PushBack(val T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.data = append(v.data, val)
}

// PopBack removes and returns the last element from the vector.
// If the vector is empty, it returns a zero value and false.
func (v *Vector[T]) PopBack() (T, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	var zero T
	if len(v.data) == 0 {
		return zero, false
	}
	idx := len(v.data) - 1
	val := v.data[idx]
	v.data = v.data[:idx]
	return val, true
}

// Len returns the number of elements in the vector.
func (v *Vector[T]) Len() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.data)
}

// Cap returns the number of elements that the vector
// can hold without further allocation.
func (v *Vector[T]) Cap() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return cap(v.data)
}

// Empty reports whether the vector is empty.
func (v *Vector[T]) Empty() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.data) == 0
}

// Get returns the zero-indexed ith element of v, if any.
func (v *Vector[T]) Get(i int) (value T, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if i >= 0 && i < len(v.data) {
		return v.data[i], true
	}
	return
}

// Set sets the zero-indexed ith element of v to value.
//
// Set panics if n is negative or greater than or equal to the length of v.
func (v *Vector[T]) Set(i int, value T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if i < 0 || i >= len(v.data) {
		panic("vector.Set: index out of range")
	}

	v.data[i] = value
}

// Front returns the first element of v, if any.
func (v *Vector[T]) Front() (value T, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if len(v.data) > 0 {
		value = v.data[0]
		ok = true
	}
	return
}

// Back returns the last element of v, if any.
func (v *Vector[T]) Back() (value T, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if len(v.data) > 0 {
		value = v.data[len(v.data)-1]
		ok = true
	}
	return
}

// Clear removes all elements from the vector
func (v *Vector[T]) Clear() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.data = v.data[:0]
}

// Grow increases the capacity of the vector, if necessary,
// to guarantee space for another n elements. After Grow(n), at least n elements
// can be appended to the vector without another allocation.
//
// If n is negative, Grow panics.
func (v *Vector[T]) Grow(n int) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if n < 0 {
		panic("vector.Grow: negative count")
	}
	if cap(v.data)-len(v.data) >= n {
		return
	}

	data := make([]T, len(v.data), len(v.data)+n)
	copy(data, v.data)
	v.data = data
}

// Clip reduces the capacity of the vector as much as possible.
func (v *Vector[T]) Clip() {
	v.mu.Lock()
	defer v.mu.Unlock()

	if len(v.data) == cap(v.data) {
		return
	}

	sz := len(v.data)
	data := make([]T, sz)
	copy(data, v.data)
	v.data = data
}

// Resize changes the size of the vector to n. If n is less than the current
// size, the vector is truncated to the first n elements. If n is greater than
// the current size, the vector is grown to n elements, with the additional
// elements initialized to zero.
//
// If n is negative, Resize panics.
func (v *Vector[T]) Resize(n int) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if n < 0 {
		panic("vector.Resize: negative count")
	}
	if n < len(v.data) {
		v.data = v.data[:n]
		return
	}
	var zero T
	for i := len(v.data); i < n; i++ {
		v.data = append(v.data, zero)
	}
}

// Insert inserts the values vals... into v at index i.
//
// Insert panics if i is out of range.
// This function is O(v.Len() + len(vals))
func (v *Vector[T]) Insert(i int, vals ...T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if i < 0 || i > len(v.data) {
		panic("vector.Insert: index out of range")
	}
	v.data = slices.Insert(v.data, i, vals...)
}

// Remove removes the element at index i from v.
//
// Remove panics if i is out of range.
func (v *Vector[T]) Remove(i int) {
	v.RemoveRange(i, i+1)
}

// RemoveRange removes the elements with indices in range [i, j) from v.
// Removing an empty range (i >= j) is a no-op.
//
// RemoveRange panics if either i or j is out of range.
func (v *Vector[T]) RemoveRange(i, j int) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if i < 0 || j < 0 || i > len(v.data) || j > len(v.data) {
		panic("vector.RemoveRange: index out of range")
	}
	if i >= j {
		return
	}
	oldData := v.data
	left := v.data[:i]
	right := v.data[j:]
	v.data = append(left, right...)
	clear(oldData[len(v.data):])
}

// Equal reports whether two vectors are equal.
func Equal[T comparable](v1, v2 *Vector[T]) bool {
	v1.mu.RLock()
	defer v1.mu.RUnlock()
	v2.mu.RLock()
	defer v2.mu.RUnlock()

	return slices.Equal(v1.data, v2.data)
}

// Values returns an iterator that yields the vector elements in order.
func (v *Vector[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		v.mu.RLock()
		defer v.mu.RUnlock()

		for _, val := range v.data {
			if !yield(val) {
				break
			}
		}
	}
}

// All returns an iterator over index-value pairs in the vector.
func (v *Vector[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		v.mu.RLock()
		defer v.mu.RUnlock()

		for i, val := range v.data {
			if !yield(i, val) {
				break
			}
		}
	}
}

// Backward returns an iterator over index-value pairs in the vector,
// traversing it backward with decreasing indices.
func (v *Vector[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		v.mu.RLock()
		defer v.mu.RUnlock()

		for i := len(v.data) - 1; i >= 0; i-- {
			if !yield(i, v.data[i]) {
				break
			}
		}
	}
}

// Collect collects values from an iterator and returns a new vector.
func Collect[T any](seq iter.Seq[T]) *Vector[T] {
	data := slices.Collect(seq)
	return &Vector[T]{
		data: data,
	}
}
