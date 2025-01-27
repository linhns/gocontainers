// Package queue implements a simple queue data structure.
package queue

// A Queue is a FIFO data structure.
type Queue[T any] struct {
	data []T
}

// New creates and initializes a new [Queue].
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Empty reports whether the queue is empty.
func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	return len(q.data)
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(val T) {
	q.data = append(q.data, val)
}

// Front returns the element at the front of the queue.
//
// It returns the zero value of T and false if the queue is empty.
// Otherwise, it returns the element and true.
func (q *Queue[T]) Front() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.data[0], true
}

// Pop removes and returns the element at the front of the queue.
//
// It returns the zero value of T and false if the queue is empty.
// Otherwise, it returns the element and true.
func (q *Queue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	val := q.data[0]
	q.data = q.data[1:]
	return val, true
}
