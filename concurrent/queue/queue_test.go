package queue_test

import (
	"sync"
	"testing"

	"github.com/linhns/gocontainers/concurrent/queue"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := queue.New[int]()
	assert.True(t, q.Empty())

	_, ok := q.Pop()
	assert.False(t, ok)

	q.Push(5)
	assert.False(t, q.Empty())

	q.Push(100)
	assert.Equal(t, 2, q.Len())

	val, ok := q.Front()
	assert.True(t, ok)
	assert.Equal(t, 5, val)

	val, _ = q.Pop()
	assert.Equal(t, 5, val)

	val, _ = q.Pop()
	assert.Equal(t, 100, val)

	assert.True(t, q.Empty())
	_, ok = q.Front()
	assert.False(t, ok)
}

func TestQueueConcurrent(t *testing.T) {
	s := queue.New[int]()

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
