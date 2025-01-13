package queue_test

import (
	"testing"

	"github.com/linhns/gocontainers/queue"
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
