package stack

import (
	"sync"

	"github.com/linhns/rwlocker"
)

type config struct {
	locker rwlocker.RWLocker
}

type Option func(*config)

func WithConcurrentUnsafe() Option {
	return func(c *config) {
		c.locker = &rwlocker.NullRWLocker{}
	}
}

type Stack[T any] struct {
	data   []T
	locker rwlocker.RWLocker
}

func New[T any](opts ...Option) *Stack[T] {
	config := config{
		locker: &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &Stack[T]{
		locker: config.locker,
	}
}

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
}

func (s *Stack[T]) Top() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

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
