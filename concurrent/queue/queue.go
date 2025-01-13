// Package queue implements a simple queue data structure
// safe for concurrent use.
package queue

import "sync"

// A Queue is a FIFO data structure.
type Queue[T any] struct {
	mu   sync.RWMutex
	data []T
}

// New creates and initializes a new [Queue].
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Empty reports whether the queue is empty.
func (q *Queue[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.data) == 0
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.data)
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(val T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.data = append(q.data, val)
}

// Front returns the element at the front of the queue.
//
// It returns the zero value of T and false if the queue is empty.
// Otherwise, it returns the element and true.
func (q *Queue[T]) Front() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

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
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	val := q.data[0]
	q.data = q.data[1:]
	return val, true
}
