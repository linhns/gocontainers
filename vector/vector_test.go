package vector_test

import (
	"math/rand/v2"
	"testing"

	"github.com/linhns/gocontainers/vector"
	"github.com/stretchr/testify/assert"
)

func TestVectorSingleElementAndBasicOperations(t *testing.T) {
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

func TestVectorCapacity(t *testing.T) {
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

func TestVectorGet(t *testing.T) {
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

func TestVectorSet(t *testing.T) {
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

func TestVectorGetSet(t *testing.T) {
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

func TestVectorClear(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	for i := 0; i < 4; i++ {
		v.Push(i)
	}

	v.Clear()
	assert.True(t, v.Empty())
	assert.Equal(t, 4, v.Capacity())
}

func TestVectorReserve(t *testing.T) {
	t.Parallel()

	v := vector.NewWithCapacity[int](4)
	v.Push(1)
	v.Push(2)
	v.Reserve(10)
	assert.True(t, v.Capacity() >= 10)
	v.Reserve(8)
	assert.True(t, v.Capacity() >= 10)
}
