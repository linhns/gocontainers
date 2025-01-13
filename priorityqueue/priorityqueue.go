package priorityqueue

import (
	"cmp"

	"github.com/linhns/gocontainers/comparator"
)

type PriorityQueue[T any] struct {
	data       []T
	comparator comparator.Comparator[T]
}

func New[T cmp.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data:       []T{},
		comparator: cmp.Compare[T],
	}
}

func NewWithComparator[T any](comparator comparator.Comparator[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data:       []T{},
		comparator: comparator,
	}
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.data)
}

func (pq *PriorityQueue[T]) Empty() bool {
	return len(pq.data) == 0
}

func (pq *PriorityQueue[T]) Push(v T) {
	pq.data = append(pq.data, v)
	pq.siftUp(len(pq.data) - 1)
}

func (pq *PriorityQueue[T]) Top() (T, bool) {
	if len(pq.data) == 0 {
		var zero T
		return zero, false
	}
	return pq.data[0], true
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
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
