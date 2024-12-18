package vector_test

import (
	"math/rand/v2"
	"testing"

	"github.com/linhns/gocontainers/vector"
	"github.com/stretchr/testify/assert"
)

func TestSingleElementAndBasicOperations(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	assert.True(t, v.Empty())
	v.Push(1)
	assert.Equal(t, 1, v.Len())
	val, ok := v.Pop()
	assert.Equal(t, 1, val)
	assert.Equal(t, true, ok)

	_, ok = v.Pop()
	assert.Equal(t, false, ok)
}

func TestCapacity(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	assert.True(t, v.Empty())
	assert.Equal(t, 4, v.Capacity())

	for i := 1; i <= 4; i++ {
		v.Push(i)
	}

	assert.Equal(t, 4, v.Capacity())
	v.Push(5)
	assert.True(t, v.Capacity() > 4)
}

func TestGet(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	for i := 0; i < 4; i++ {
		v.Push(i)
	}

	tests := map[string]struct {
		index     int
		wantValue int
		wantOk    bool
	}{
		"valid index":     {index: 2, wantValue: 2, wantOk: true},
		"negative index":  {index: -1, wantOk: false},
		"too large index": {index: 4, wantOk: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			val, ok := v.Get(tt.index)
			assert.Equal(t, tt.wantValue, val)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestSet(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	for i := 0; i < 4; i++ {
		v.Push(i)
	}

	tests := map[string]struct {
		index  int
		value  int
		wantOk bool
	}{
		"valid index":     {index: 2, value: 2, wantOk: true},
		"negative index":  {index: -1, wantOk: false},
		"too large index": {index: 4, wantOk: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ok := v.Set(tt.index, tt.value)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestGetSet(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	for i := 0; i < 4; i++ {
		v.Push(i)
	}

	for i := 0; i < 4; i++ {
		want := rand.Int()
		_ = v.Set(i, want)
		got, _ := v.Get(i)
		assert.Equal(t, want, got)
	}
}

func TestFront(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	_, ok := v.Front()
	assert.False(t, ok)
	v.Push(1)
	v.Push(2)
	v.Push(3)
	val, ok := v.Front()
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestBack(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	_, ok := v.Back()
	assert.False(t, ok)
	v.Push(1)
	v.Push(2)
	v.Push(3)
	val, ok := v.Back()
	assert.True(t, ok)
	assert.Equal(t, 3, val)
}

func TestClear(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	for i := 0; i < 4; i++ {
		v.Push(i)
	}

	v.Clear()
	assert.True(t, v.Empty())
	assert.Equal(t, 4, v.Capacity())
}

func TestReserve(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	v.Push(1)
	v.Push(2)
	v.Reserve(10)
	assert.True(t, v.Capacity() >= 10)
	v.Reserve(8)
	assert.True(t, v.Capacity() >= 10)
}

func TestShrinkToFit(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](10)
	v.Push(1)
	v.Push(2)
	v.Push(3)
	v.ShrinkToFit()
	assert.Equal(t, 3, v.Capacity())
}

func TestResize(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](1024)
	for i := 0; i < 1024; i++ {
		v.Push(1)
	}

	v.Resize(1000)
	assert.Equal(t, 1000, v.Len())

	v.Resize(1024)

	for i := 1000; i < 1024; i++ {
		val, ok := v.Get(i)
		assert.Equal(t, 0, val)
		assert.True(t, ok)
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	v.Push(1)
	v.Push(2)
	v.Push(3)
	v.Insert(1, 4)
	assert.Equal(t, 4, v.Len())
	val, _ := v.Get(1)
	assert.Equal(t, 4, val)

	v.Insert(4, 5)
}

func TestEqual(t *testing.T) {
	t.Parallel()

	v1 := vector.New[int]()
	v2 := vector.New[int]()
	assert.True(t, vector.Equal(v1, v2))

	v1.Push(1)
	assert.False(t, vector.Equal(v1, v2))

	v2.Push(1)
	assert.True(t, vector.Equal(v1, v2))

	v1.Push(2)
	v2.Push(3)
	assert.False(t, vector.Equal(v1, v2))
}

func TestCollectValues(t *testing.T) {
	t.Parallel()

	v := vector.New[int]()
	for i := 0; i < 10; i++ {
		v.Push(i)
	}

	assert.True(t, vector.Equal(v, vector.Collect(v.Values())))
}

func TestAll(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](5)
	elems := []int{1, 2, 3, 4, 5}
	for _, elem := range elems {
		v.Push(elem)
	}

	idx := 0
	for i, val := range v.All() {
		assert.Equal(t, i, idx)
		assert.Equal(t, elems[i], val)
		idx++
	}
}
