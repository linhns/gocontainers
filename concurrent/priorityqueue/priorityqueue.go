// Package priorityqueue provides a generic priority queue implementation
// based on array heap that is safe for concurrent use.
package priorityqueue

import (
	"sync"

	"github.com/linhns/gocontainers/comparator"
)

// PriorityQueue is a generic priority queue with a configurable comparision
// function (comparator).
type PriorityQueue[T any] struct {
	mu         sync.RWMutex
	data       []T
	comparator comparator.Comparator[T]
}

// New creates a new [PriorityQueue] with the specified comparator.
func New[T any](comparator comparator.Comparator[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data:       []T{},
		comparator: comparator,
	}
}

// Len returns the number of elements in the priority queue.
func (pq *PriorityQueue[T]) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.data)
}

// Empty reports whether the priority queue is empty.
func (pq *PriorityQueue[T]) Empty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.data) == 0
}

// Push adds an element to the priority queue.
func (pq *PriorityQueue[T]) Push(v T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pq.data = append(pq.data, v)
	pq.siftUp(len(pq.data) - 1)
}

// Top returns the element with the maximal priority in the queue.
// If the queue is empty, it returns the zero value of the element type
// and false.
func (pq *PriorityQueue[T]) Top() (T, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.data) == 0 {
		var zero T
		return zero, false
	}
	return pq.data[0], true
}

// Pop returns the element with the maximal priority in the queue, and
// removes it from the queue. If the queue is empty, it returns the zero
// value of the element type and false.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.data) == 0 {
		var zero T
		return zero, false
	}
	top := pq.data[0]
	pq.data[0] = pq.data[len(pq.data)-1]
	pq.data = pq.data[:len(pq.data)-1]
	pq.siftDown(0)
	return top, true
}

func (pq *PriorityQueue[T]) siftUp(index int) {
	for {
		parent := (index - 1) / 2
		if index <= 0 || pq.comparator(pq.data[index], pq.data[parent]) <= 0 {
			break
		}
		pq.data[parent], pq.data[index] = pq.data[index], pq.data[parent]
		index = parent
	}
}

func (pq *PriorityQueue[T]) siftDown(index int) {
	i := index
	left := 2*i + 1
	right := 2*i + 2

	if left < len(pq.data) && pq.comparator(pq.data[left], pq.data[i]) > 0 {
		i = left
	}

	if right < len(pq.data) && pq.comparator(pq.data[right], pq.data[i]) > 0 {
		i = right
	}

	if i != index {
		pq.data[index], pq.data[i] = pq.data[i], pq.data[index]
		pq.siftDown(i)
	}
}
