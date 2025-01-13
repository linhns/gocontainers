package comparator_test

import (
	"cmp"
	"testing"

	"github.com/linhns/gocontainers/comparator"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	cmpr := comparator.Reverse(cmp.Compare[int])

	assert.Equal(t, 0, cmpr(5, 5))
	assert.Equal(t, 1, cmpr(5, 6))
	assert.Equal(t, -1, cmpr(6, 5))
}
