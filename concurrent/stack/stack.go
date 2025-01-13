package stack

import (
	"sync"
)

type Stack[T any] struct {
	mu   sync.RWMutex
	data []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Empty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Stack[T]) Push(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, val)
}

func (s *Stack[T]) Top() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	index := len(s.data) - 1
	val := s.data[index]
	s.data = s.data[:index]
	return val, true
}
