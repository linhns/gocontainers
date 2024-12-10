package set_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"

	"github.com/linhns/gocontainers/set"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOfDistinct(rapid.IntMin(0), rapid.ID)
		data := gen.Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}
		assert.Equal(t, len(data), s.Len())

		if len(data) > 0 {
			existingData := rapid.SliceOf(rapid.SampledFrom(data)).Draw(t, "existing")
			for _, v := range existingData {
				s.Add(v)
			}
			assert.Equal(t, len(data), s.Len(), "adding existing elements should not change size")
		}

		newData := rapid.SliceOfDistinct(rapid.IntMax(-1), rapid.ID).Draw(t, "new")
		for _, v := range newData {
			s.Add(v)
		}
		assert.Equal(t, len(data)+len(newData), s.Len(), "adding new elements should change size by the number of new elements")
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()
	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOf(rapid.Int()).Draw(t, "data")
		s1 := set.New[int]()
		s2 := set.New[int]()
		for _, v := range data {
			s1.Add(v)
			s2.Add(v)
		}
		assert.True(t, set.Equal(s1, s2))
	})
}

func TestContains(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOf(rapid.IntMin(0)).Draw(t, "data")
		s := set.New[int]()

		for _, v := range data {
			s.Add(v)
		}

		for _, v := range data {
			assert.True(t, s.Contains(v))
		}

		notInData := rapid.SliceOf(rapid.IntMax(-1)).Draw(t, "notInData")
		for _, v := range notInData {
			assert.False(t, s.Contains(v))
		}
	})
}

func TestLen(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOfDistinct(rapid.Int(), rapid.ID).Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}

		assert.Equal(t, len(data), s.Len())
	})
}

func TestClear(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOf(rapid.Int()).Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}

		s.Clear()
		assert.Equal(t, 0, s.Len())
	})
}

func TestRemove(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOfDistinct(rapid.IntMin(0), rapid.ID)
		data := gen.Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}

		if len(data) > 0 {
			existing := rapid.SliceOfDistinct(rapid.SampledFrom(data), rapid.ID).Draw(t, "existing")
			for _, v := range existing {
				oldSize := s.Len()
				s.Remove(v)
				assert.False(t, s.Contains(v))
				assert.Equal(t, oldSize-1, s.Len(), "removing an existing elements should decrease size by 1")
			}
		}

		oldSize := s.Len()
		notExists := rapid.SliceOfDistinct(rapid.IntMax(-1), rapid.ID).Draw(t, "notExists")
		for _, v := range notExists {
			s.Remove(v)
		}
		assert.Equal(t, oldSize, s.Len(), "removing non-existing elements should not change size")
	})
}

func TestAll(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOfDistinct(rapid.Int(), rapid.ID).Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}

		assert.True(t, slices.Equal(slices.Sorted(slices.Values(data)), slices.Sorted(s.All())))
	})
}

func TestCollect(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOfDistinct(rapid.Int(), rapid.ID).Draw(t, "data")
		seq := slices.Values(data)
		s1 := set.Collect(seq)
		s2 := set.New[int]()
		for _, v := range data {
			s2.Add(v)
		}
		assert.True(t, set.Equal(s1, s2))
	})
}

func TestAllCollectRoundtrip(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		data := rapid.SliceOfDistinct(rapid.Int(), rapid.ID).Draw(t, "data")
		s := set.New[int]()
		for _, v := range data {
			s.Add(v)
		}
		rdtrip := set.Collect(s.All())
		assert.True(t, set.Equal(s, rdtrip))
	})
}
