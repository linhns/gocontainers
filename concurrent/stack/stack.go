// Package stack provides a simple stack implementation
// safe for concurrent use.
package stack

import (
	"sync"
)

// A Stack is a Last-In-First-Out (LIFO) data structure.
type Stack[T any] struct {
	mu   sync.RWMutex
	data []T
}

// New creates and initializes a new [Stack].
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Empty reports whether the stack is empty.
func (s *Stack[T]) Empty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data) == 0
}

// Len returns the number of elements in the stack.
func (s *Stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, val)
}

// Top returns the element at the top of the stack.
//
// If the stack is empty, it returns the zero value of T and false.
func (s *Stack[T]) Top() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Pop removes and returns the element at the top of the stack.
//
// If the stack is empty, it returns the zero value of T and false.
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
