package hashset_test

import (
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linhns/gocontainers/concurrent/hashset"
)

func TestSetBasicOperations(t *testing.T) {
	t.Parallel()

	s := hashset.New[int]()

	assert.True(t, s.Empty())
	s.Add(1)
	assert.False(t, s.Empty())

	assert.False(t, s.Contains(10))
	s.Add(10)
	assert.True(t, s.Contains(10))

	assert.Equal(t, 2, s.Len())
	s.Add(1)
	assert.Equal(t, 2, s.Len())

	s.Remove(10)
	assert.False(t, s.Contains(10))

	assert.Equal(t, 1, s.Len())
	s.Remove(4)
	assert.Equal(t, 1, s.Len())

	s.Clear()
	assert.True(t, s.Empty())
}

func TestSetEqual(t *testing.T) {
	t.Parallel()

	s1 := hashset.New[int]()
	s2 := hashset.New[int]()

	assert.True(t, hashset.Equal(s1, s2))
	s1.Add(1)
	s1.Add(2)
	assert.False(t, hashset.Equal(s1, s2))

	s2.Add(2)
	s2.Add(1)
	assert.True(t, hashset.Equal(s1, s2))

	s3 := hashset.New[int]()
	s3.Add(1)
	s3.Add(3)
	assert.False(t, hashset.Equal(s1, s3))
}

func TestSetAllCollectRoundtrip(t *testing.T) {
	t.Parallel()

	s := hashset.New[string]()
	s.Add("dog")
	s.Add("cat")
	s.Add("fox")
	s.Add("fox")

	rdtrip := hashset.Collect(s.All())

	assert.True(t, hashset.Equal(s, rdtrip))
}

func TestSetMathOperations(t *testing.T) {
	t.Parallel()

	s1 := hashset.New[int]()
	s2 := hashset.New[int]()

	s1.Add(1)
	s1.Add(2)
	s1.Add(4)
	s1.Add(5)

	s2.Add(1)
	s2.Add(2)
	s2.Add(3)
	s2.Add(4)

	union := hashset.Collect(slices.Values([]int{1, 2, 3, 4, 5}))
	assert.True(t, hashset.Equal(hashset.Union(s1, s2), union))

	intersection := hashset.Collect(slices.Values([]int{1, 2, 4}))
	assert.True(t, hashset.Equal(hashset.Intersection(s1, s2), intersection))

	difference := hashset.Collect(slices.Values([]int{5}))
	assert.True(t, hashset.Equal(hashset.Difference(s1, s2), difference))
}

func TestSetConcurrent(t *testing.T) {
	s := hashset.New[int]()

	var wg sync.WaitGroup
	wg.Add(200)

	start := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			<-start
			s.Add(i)
		}(i)

		go func(i int) {
			defer wg.Done()
			<-start
			s.Remove(i)
		}(i)
	}

	close(start)
	wg.Wait()
}
