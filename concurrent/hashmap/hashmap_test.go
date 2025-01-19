package hashmap_test

import (
	"maps"
	"slices"
	"sync"
	"testing"

	"github.com/linhns/gocontainers/concurrent/hashmap"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Parallel()

	m := hashmap.New[string, int]()

	assert.True(t, m.Empty())
	assert.False(t, m.Contains("one"))

	m.Insert("one", 1)
	assert.True(t, m.Contains("one"))
	assert.Equal(t, 1, m.Len())

	m.Remove("one")
	assert.Equal(t, 0, m.Len())

	m.Insert("two", 2)
	m.Insert("two", 2)
	m.Insert("three", 3)
	assert.Equal(t, 2, m.Len())

	m.Remove("four")
	assert.Equal(t, 2, m.Len())

	_, ok := m.Get("four")
	assert.False(t, ok)

	val, ok := m.Get("three")
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	m.Clear()
	assert.True(t, m.Empty())
}

func TestMapIterator(t *testing.T) {
	t.Parallel()

	m := hashmap.New[int, string]()

	ints := []int{1, 2, 3, 4, 5}
	strs := []string{"one", "two", "three", "four", "five"}

	for i := range ints {
		m.Insert(ints[i], strs[i])
	}

	assert.ElementsMatch(t, ints, slices.Collect(m.Keys()))
	assert.ElementsMatch(t, strs, slices.Collect(m.Values()))

	want := maps.Collect(m.All())
	got := maps.Collect(hashmap.Collect(m.All()).All())
	assert.Equal(t, want, got)
}

func TestMapConcurrent(t *testing.T) {
	t.Parallel()

	m := hashmap.New[string, int]()

	ints := []int{1, 2, 3, 4, 5}
	strs := []string{"one", "two", "three", "four", "five"}

	var wg sync.WaitGroup
	start := make(chan struct{})

	addFunc := func(key string, val int) {
		defer wg.Done()
		<-start
		m.Insert(key, val)
	}

	removeFunc := func(key string) {
		defer wg.Done()
		<-start
		m.Remove(key)
	}

	wg.Add(2 * len(ints))
	for i := range ints {
		go addFunc(strs[i], ints[i])
	}

	for i := range ints {
		go removeFunc(strs[i])
	}

	close(start)
	wg.Wait()

	for k, v := range m.All() {
		assert.True(t, slices.Contains(strs, k))
		assert.True(t, slices.Contains(ints, v))
	}
}

func TestMapIteratorConcurrent(t *testing.T) {
	t.Parallel()

	m := hashmap.New[string, int]()

	ints := []int{1, 2, 3, 4, 5}
	strs := []string{"one", "two", "three", "four", "five"}

	for i := range ints {
		m.Insert(strs[i], ints[i])
	}

	start := make(chan struct{})

	go func() {
		<-start
		for range m.Keys() {
		}
	}()

	go func() {
		<-start
		for range m.Values() {
		}
	}()

	go func() {
		<-start
		hashmap.Collect(m.All())
	}()
	close(start)
}
