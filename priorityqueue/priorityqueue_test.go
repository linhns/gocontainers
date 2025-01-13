package priorityqueue_test

import (
	"cmp"
	"testing"

	"github.com/linhns/gocontainers/priorityqueue"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueueCustomType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	cmpr := func(a, b Person) int {
		if a.Age == b.Age {
			return cmp.Compare(a.Name, b.Name)
		}
		return cmp.Compare(a.Age, b.Age)
	}
	pq := priorityqueue.NewWithComparator(cmpr)

	assert.True(t, pq.Empty())

	_, ok := pq.Pop()
	assert.False(t, ok)

	pq.Push(Person{"Danny", 30})
	pq.Push(Person{"John", 30})
	pq.Push(Person{"Adam", 25})

	assert.Equal(t, 3, pq.Len())

	top, ok := pq.Top()
	assert.True(t, ok)
	assert.Equal(t, Person{"John", 30}, top)

	pq.Pop()
	v, ok := pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, Person{"Danny", 30}, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, Person{"Adam", 25}, v)

	assert.True(t, pq.Empty())

	_, ok = pq.Top()
	assert.False(t, ok)
}
