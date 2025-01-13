// Package stack provides a simple stack implementation.
package stack

// A Stack is a Last-In-First-Out (LIFO) data structure.
type Stack[T any] struct {
	data []T
}

// New creates and initializes a new [Stack].
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Empty reports whether the stack is empty.
func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

// Len returns the number of elements in the stack.
func (s *Stack[T]) Len() int {
	return len(s.data)
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
}

// Top returns the element at the top of the stack.
//
// If the stack is empty, it returns the zero value of T and false.
func (s *Stack[T]) Top() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Pop removes and returns the element at the top of the stack.
//
// If the stack is empty, it returns the zero value of T and false.
func (s *Stack[T]) Pop() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	index := len(s.data) - 1
	val := s.data[index]
	s.data = s.data[:index]
	return val, true
}
