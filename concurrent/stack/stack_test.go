package stack_test

import (
	"sync"
	"testing"

	"github.com/linhns/gocontainers/concurrent/stack"
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

func TestStackConcurrent(t *testing.T) {
	s := stack.New[int]()

	var wg sync.WaitGroup
	wg.Add(200)

	start := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			<-start
			s.Push(i)
		}(i)

		go func() {
			defer wg.Done()
			<-start
			_, _ = s.Pop()
		}()
	}

	close(start)
	wg.Wait()
}
