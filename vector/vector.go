// Package vector provides a simple vector implementation
// based on Go slices.
package vector

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

// Push adds an element to the end of the collection
func (v *Vector[T]) Push(val T) {
	v.data = append(v.data, val)
}

// Pop removes and returns the last element from the collection.
// If the collection is empty, it returns a zero value and false.
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

// Len returns the number of elements in the collection
func (v *Vector[T]) Len() int {
	return len(v.data)
}

// Capacity returns the number of elements that the collection
// can hold without further allocation
func (v *Vector[T]) Capacity() int {
	return cap(v.data)
}

// Empty reports whether the collection is empty
func (v *Vector[T]) Empty() bool {
	return len(v.data) == 0
}

// Get returns the element at the specified index.
// If the index is out of range, it returns a zero value and false.
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

// Clear removes all elements from the collection
func (v *Vector[T]) Clear() {
	v.data = v.data[:0]
}

// Reserve increases the capacity of the collection to
// the specified value.
//
// If the specified capacity is less than or equal to the current capacity,
// this is a no-op.
func (v *Vector[T]) Reserve(capacity int) {
	if capacity <= cap(v.data) {
		return
	}
	newData := make([]T, len(v.data), capacity)
	copy(newData, v.data)
	v.data = newData
}
