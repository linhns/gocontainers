package priorityqueue_test

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"sync"
	"testing"

	"github.com/linhns/gocontainers/concurrent/priorityqueue"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueueInt(t *testing.T) {
	pq := priorityqueue.New(cmp.Compare[int])

	ints := []int{10, 20, 30, 40, 50}
	for i := 0; i < 100; i++ {
		nums := slices.Clone(ints)
		rand.Shuffle(len(ints), func(i, j int) {
			nums[i], nums[j] = nums[j], nums[i]
		})

		for _, num := range nums {
			pq.Push(num)
		}

		val, _ := pq.Pop()
		assert.Equal(t, 50, val)
		assert.Equal(t, 4, pq.Len())

		val, _ = pq.Pop()
		assert.Equal(t, 40, val)
		assert.Equal(t, 3, pq.Len())

		val, _ = pq.Pop()
		assert.Equal(t, 30, val)
		assert.Equal(t, 2, pq.Len())

		val, _ = pq.Pop()
		assert.Equal(t, 20, val)
		assert.Equal(t, 1, pq.Len())

		val, _ = pq.Pop()
		assert.Equal(t, 10, val)
		assert.Equal(t, 0, pq.Len())
	}
}

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
	pq := priorityqueue.New(cmpr)

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

func TestPriorityQueueConcurrent(t *testing.T) {
	pq := priorityqueue.New(cmp.Compare[string])

	var wg sync.WaitGroup
	wg.Add(200)

	start := make(chan struct{})
	words := []string{"apple", "banana", "cherry", "date"}
	for i := 0; i < 100; i++ {
		word := words[i%4]
		go func(s string) {
			defer wg.Done()
			<-start
			pq.Push(s)
		}(word)

		go func() {
			defer wg.Done()
			<-start
			_, _ = pq.Pop()
		}()
	}

	close(start)
	wg.Wait()
}
