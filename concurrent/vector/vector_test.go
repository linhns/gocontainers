package vector_test

import (
	"sync"
	"testing"

	"github.com/linhns/gocontainers/concurrent/vector"
	"github.com/stretchr/testify/assert"
)

func TestVectorBasics(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.True(t, v.Empty())
	assert.Equal(t, 0, v.Cap())

	val, ok := v.Front()
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	val, ok = v.Back()
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	v.PushBack(1)
	assert.Equal(t, 1, v.Len())

	val, ok = v.Front()
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	for i := 2; i <= 5; i++ {
		v.PushBack(i)
		val, ok := v.Back()
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}
	assert.Equal(t, 5, v.Len())

	val, ok = v.Front()
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = v.PopBack()
	assert.True(t, ok)
	assert.Equal(t, 5, val)

	assert.Equal(t, 4, v.Len())

	val, ok = v.Get(1000)
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	val, ok = v.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	v.Set(2, 1000)
	val, _ = v.Get(2)
	assert.Equal(t, 1000, val)

	v.Set(2, 100)
	v.Set(2, 200)
	val, _ = v.Get(2)
	assert.Equal(t, 200, val)

	v.Clear()
	assert.True(t, v.Empty())

	val, ok = v.PopBack()
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	v.PushBack(1)
	v.PushBack(2)
	v.Insert(1, 1, 2)

	val, ok = v.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	v.RemoveRange(1, 3)
	assert.Equal(t, 2, v.Len())

	v.Remove(0)
	val, ok = v.Front()
	assert.True(t, ok)
	assert.Equal(t, 2, val)
}

func TestVectorCapacity(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	assert.True(t, v.Empty())
	assert.Equal(t, 4, v.Cap())

	for i := 1; i <= 4; i++ {
		v.PushBack(i)
	}

	assert.Equal(t, 4, v.Cap())
	v.PushBack(5)
	assert.True(t, v.Cap() > 4)
}

func TestVectorOf(t *testing.T) {
	t.Parallel()

	v := vector.Of(1, 2, 3, 4)
	assert.Equal(t, 4, v.Len())
	assert.Equal(t, 4, v.Cap())

	for i := 0; i < 4; i++ {
		val, ok := v.Get(i)
		assert.True(t, ok)
		assert.Equal(t, i+1, val)
	}
}

func TestVectorGrowNegative(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Grow(-1) })
}

func TestVectorGrowNoOp(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](10)
	v.Grow(4)
	assert.Equal(t, 10, v.Cap())
}

func TestVectorGrow(t *testing.T) {
	t.Parallel()

	v := vector.Of(1, 2)
	v.Grow(5)
	assert.GreaterOrEqual(t, v.Cap()-v.Len(), 5)
}

func TestVectorClip(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](5)
	v.PushBack(1)
	v.PushBack(2)
	v.Clip()

	assert.Equal(t, 2, v.Len())
	assert.Equal(t, 2, v.Cap())

	v.Clip()
	assert.Equal(t, 2, v.Cap())
}

func TestVectorResizeNegative(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Resize(-1) })
}

func TestVectorResize(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	v.PushBack(1)
	v.PushBack(2)
	v.Resize(4)

	val, ok := v.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	val, ok = v.Get(3)
	assert.True(t, ok)
	assert.Equal(t, 0, val)

	v.Resize(1)
	assert.Equal(t, 1, v.Len())

	val, ok = v.Get(0)
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestVectorSetIndexNegative(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Set(-1, 0) })
}

func TestVectorSetIndexTooLarge(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Set(0, 0) })
}

func TestVectorInsertIndexNegative(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Insert(-1, 0, 1) })
}

func TestVectorInsertIndexTooLarge(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	v.PushBack(1)
	assert.Panics(t, func() { v.Insert(2, 0, 1) })
}

func TestVectorInsert(t *testing.T) {
	t.Parallel()

	initialVals := []int{11, 12, 13, 14, 15}

	tests := map[string]struct {
		index int
		vals  []int
	}{
		"at front": {
			index: 0,
			vals:  []int{1, 2, 3, 4, 5},
		},
		"at back": {
			index: 5,
			vals:  []int{1, 2, 3, 4, 5},
		},
		"in middle": {
			index: 3,
			vals:  []int{1, 2, 3, 4, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := vector.Of(initialVals...)

			v.Insert(tt.index, tt.vals...)
			for i := 0; i < tt.index; i++ {
				val, _ := v.Get(i)
				assert.Equal(t, initialVals[i], val)
			}
			for i := tt.index; i < tt.index+len(tt.vals); i++ {
				val, _ := v.Get(i)
				assert.Equal(t, tt.vals[i-tt.index], val)
			}
			for i := tt.index + len(tt.vals); i < len(initialVals)+len(tt.vals); i++ {
				val, _ := v.Get(i)
				assert.Equal(t, initialVals[i-len(tt.vals)], val)
			}
		})
	}
}

func TestVectorRemoveNegative(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Remove(-1) })
}

func TestVectorRemoveIndexTooLarge(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.Panics(t, func() { v.Remove(0) })
}

func TestVectorRemoveRange(t *testing.T) {
	t.Parallel()

	v := vector.Of(1, 2, 3, 4, 5)

	v.RemoveRange(5, 5)
	assert.Equal(t, 5, v.Len())

	v.RemoveRange(1, 4)
	val, ok := v.Get(0)
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = v.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

func TestVectorEqual(t *testing.T) {
	t.Parallel()

	v1 := vector.New[int]()
	v2 := vector.New[int]()
	assert.True(t, vector.Equal(v1, v2))

	v1.PushBack(1)
	assert.False(t, vector.Equal(v1, v2))

	v2.PushBack(1)
	assert.True(t, vector.Equal(v1, v2))

	v1.PushBack(2)
	v2.PushBack(3)
	assert.False(t, vector.Equal(v1, v2))
}

func TestCollect(t *testing.T) {
	testSeq := func(yield func(int) bool) {
		for i := 0; i < 10; i += 2 {
			if !yield(i) {
				return
			}
		}
	}
	want := vector.Of(0, 2, 4, 6, 8)

	got := vector.Collect(testSeq)
	assert.True(t, vector.Equal(got, want))
}

func TestAll(t *testing.T) {
	for size := 0; size < 10; size++ {
		v := vector.New[int]()
		for i := range size {
			v.PushBack(i)
		}
		ei, ev := 0, 0
		cnt := 0
		for i, num := range v.All() {
			assert.Equal(t, ei, i)
			assert.Equal(t, ev, num)
			ei++
			ev++
			cnt++
		}
		assert.Equal(t, size, cnt)
	}
}

func TestValues(t *testing.T) {
	for size := 0; size < 10; size++ {
		v := vector.New[int]()
		for i := range size {
			v.PushBack(i)
		}
		ev := 0
		cnt := 0
		for num := range v.All() {
			assert.Equal(t, ev, num)
			ev++
			cnt++
		}
		assert.Equal(t, size, cnt)
	}
}

func TestBackward(t *testing.T) {
	for size := 0; size < 10; size++ {
		v := vector.New[int]()
		for i := range size {
			v.PushBack(i)
		}
		ei, ev := size-1, size-1
		cnt := 0
		for i, num := range v.Backward() {
			assert.Equal(t, ei, i)
			assert.Equal(t, ev, num)
			ei--
			ev--
			cnt++
		}
		assert.Equal(t, size, cnt)
	}
}

func TestValuesCollectRoundtrip(t *testing.T) {
	want := vector.Of(0, 2, 4, 6, 8)
	got := vector.Collect(want.Values())
	assert.True(t, vector.Equal(got, want))
}

func TestVectorConcurrent(t *testing.T) {
	t.Parallel()

	strs := []string{"one", "two", "three", "four", "five"}
	v := vector.Of(strs...)

	var wg sync.WaitGroup
	start := make(chan struct{})

	addFunc := func(strs []string) {
		defer wg.Done()
		<-start
		v.Insert(1, strs...)
	}

	removeFunc := func() {
		defer wg.Done()
		<-start
		v.RemoveRange(0, 4)
	}

	pushFunc := func() {
		defer wg.Done()
		<-start
		v.PushBack("six")
	}

	popFunc := func() {
		defer wg.Done()
		<-start
		_, _ = v.PopBack()
	}

	wg.Add(4)
	go addFunc(strs)
	go removeFunc()
	go pushFunc()
	go popFunc()

	close(start)
	wg.Wait()
}

func TestVectorIteratorConcurrent(t *testing.T) {
	t.Parallel()

	v := vector.Of(1, 2, 3, 4, 5)

	start := make(chan struct{})

	go func() {
		<-start
		for range v.Backward() {
		}
	}()

	go func() {
		<-start
		for range v.All() {
		}
	}()

	go func() {
		<-start
		vector.Collect(v.Values())
	}()
	close(start)
}
