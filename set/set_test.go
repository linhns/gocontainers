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
		data := rapid.SliceOfNDistinct(rapid.Int(), 5, 5, rapid.ID).Draw(t, "data")
		s1 := set.New[int]()
		s2 := set.New[int]()
		s1.Add(data[0])
		assert.False(t, set.Equal(s1, s2), "sets of different sizes are not equal")
		s2.Add(data[1])
		assert.False(t, set.Equal(s1, s2), "sets with different elements are not equal")
		for _, v := range data {
			s1.Add(v)
			s2.Add(v)
		}
		assert.True(t, set.Equal(s1, s2), "sets with the same elements are equal")
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

		dummy := make([]int, 0)
		f := func(x int) bool {
			dummy = append(dummy, x)
			return false
		}
		s.All()(f)

		assert.Equal(t, min(len(data), 1), len(dummy))
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

func TestUnion(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		s1 := set.New[int]()
		s2 := set.New[int]()
		gen := rapid.SliceOf(rapid.Int())
		data1 := gen.Draw(t, "data1")
		data2 := gen.Draw(t, "data2")
		for _, v := range data1 {
			s1.Add(v)
		}
		for _, v := range data2 {
			s2.Add(v)
		}

		u := set.Union(s1, s2)
		expected := slices.Concat(data1, data2)
		slices.Sort(expected)
		expected = slices.Compact(expected)
		got := slices.Collect(u.All())
		slices.Sort(got)
		assert.Equal(t, expected, got)

		if len(data1)+len(data2) > 0 {
			s3 := set.New[int]()
			dups := rapid.SliceOf(rapid.SampledFrom(expected)).Draw(t, "dups")
			for _, v := range dups {
				s3.Add(v)
			}
			u2 := set.Union(u, s3)
			got = slices.Collect(u2.All())
			slices.Sort(got)
			assert.Equal(t, expected, got)
		}
	})
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		s1 := set.New[int]()
		s2 := set.New[int]()
		gen := rapid.SliceOfDistinct(rapid.Int(), rapid.ID)
		data1 := gen.Draw(t, "data1")
		data2 := gen.Draw(t, "data2")
		for _, v := range data1 {
			s1.Add(v)
		}
		for _, v := range data2 {
			s2.Add(v)
		}

		i := set.Intersection(s1, s2)

		var expected []int
		for _, v := range data1 {
			if slices.Contains(data2, v) {
				expected = append(expected, v)
			}
		}
		slices.Sort(expected)
		got := slices.Collect(i.All())
		slices.Sort(got)
		assert.Equal(t, expected, got)
	})
}

func TestDifference(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		s1 := set.New[int]()
		s2 := set.New[int]()
		gen := rapid.SliceOfDistinct(rapid.Int(), rapid.ID)
		data1 := gen.Draw(t, "data1")
		data2 := gen.Draw(t, "data2")
		for _, v := range data1 {
			s1.Add(v)
		}
		for _, v := range data2 {
			s2.Add(v)
		}

		d := set.Difference(s1, s2)

		var expected []int
		for _, v := range data1 {
			if !slices.Contains(data2, v) {
				expected = append(expected, v)
			}
		}
		slices.Sort(expected)
		got := slices.Collect(d.All())
		slices.Sort(got)
		assert.Equal(t, expected, got)
	})
}

func TestSetIdempotence(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		s := set.New[int]()
		for _, v := range rapid.SliceOf(rapid.Int()).Draw(t, "data") {
			s.Add(v)
		}
		assert.True(t, set.Equal(s, set.Union(s, s)))
		assert.True(t, set.Equal(s, set.Intersection(s, s)))
		assert.True(t, set.Equal(set.New[int](), set.Difference(s, s)))
	})
}

func TestSetCommutativity(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOf(rapid.Int())
		a := set.New[int]()
		b := set.New[int]()
		for _, v := range gen.Draw(t, "data1") {
			a.Add(v)
		}
		for _, v := range gen.Draw(t, "data2") {
			b.Add(v)
		}
		assert.True(t, set.Equal(set.Intersection(a, b), set.Intersection(b, a)))
		assert.True(t, set.Equal(set.Union(a, b), set.Union(b, a)))
	})
}

func TestSetAssociativity(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOf(rapid.Int())
		a := set.New[int]()
		b := set.New[int]()
		c := set.New[int]()
		for _, v := range gen.Draw(t, "data1") {
			a.Add(v)
		}
		for _, v := range gen.Draw(t, "data2") {
			b.Add(v)
		}
		for _, v := range gen.Draw(t, "data3") {
			c.Add(v)
		}
		assert.True(t, set.Equal(set.Intersection(a, set.Intersection(b, c)), set.Intersection(set.Intersection(a, b), c)))
		assert.True(t, set.Equal(set.Union(a, set.Union(b, c)), set.Union(set.Union(a, b), c)))
	})
}

func TestSetDistributivity(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOf(rapid.Int())
		a := set.New[int]()
		b := set.New[int]()
		c := set.New[int]()
		for _, v := range gen.Draw(t, "data1") {
			a.Add(v)
		}
		for _, v := range gen.Draw(t, "data2") {
			b.Add(v)
		}
		for _, v := range gen.Draw(t, "data3") {
			c.Add(v)
		}
		assert.True(t, set.Equal(set.Intersection(a, set.Union(b, c)), set.Union(set.Intersection(a, b), set.Intersection(a, c))))
		assert.True(t, set.Equal(set.Union(a, set.Intersection(b, c)), set.Intersection(set.Union(a, b), set.Union(a, c))))
	})
}

func TestSetAbsorption(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		gen := rapid.SliceOf(rapid.Int())
		a := set.New[int]()
		b := set.New[int]()
		for _, v := range gen.Draw(t, "data1") {
			a.Add(v)
		}
		for _, v := range gen.Draw(t, "data2") {
			b.Add(v)
		}

		assert.True(t, set.Equal(set.Intersection(a, set.Union(a, b)), a))
		assert.True(t, set.Equal(set.Union(a, set.Intersection(a, b)), a))
	})
}
