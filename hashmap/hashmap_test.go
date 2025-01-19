package hashmap_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/linhns/gocontainers/hashmap"
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
