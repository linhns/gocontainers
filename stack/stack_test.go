package stack_test

import (
	"testing"

	"github.com/linhns/gocontainers/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := stack.New[int]()
	assert.True(t, s.Empty())

	s.Push(5)
	assert.False(t, s.Empty())

	s.Push(100)
	assert.Equal(t, 2, s.Len())

	val, _ := s.Top()
	assert.Equal(t, 100, val)

	val, _ = s.Pop()
	assert.Equal(t, 100, val)

	val, _ = s.Pop()
	assert.Equal(t, 5, val)

	assert.True(t, s.Empty())
	_, ok := s.Top()
	assert.False(t, ok)
}
