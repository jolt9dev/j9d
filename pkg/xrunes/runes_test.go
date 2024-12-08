package xrunes_test

import (
	"testing"

	runes "github.com/jolt9dev/j9d/pkg/xrunes"
	"github.com/stretchr/testify/assert"
)

func TestEqualFold(t *testing.T) {
	assert.Equal(t, runes.EqualFold([]rune("test"), []rune("TEST")), true)
	assert.Equal(t, runes.EqualFold([]rune("test"), []rune("Test")), true)
	assert.Equal(t, runes.EqualFold([]rune("test"), []rune("tEsT")), true)
	assert.Equal(t, runes.EqualFold([]rune("test"), []rune("teSt")), true)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("test")), true)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("test ")), false)
	assert.Equal(t, runes.Equal([]rune("test"), []rune(" test")), false)
}

func TestEqual(t *testing.T) {
	assert.Equal(t, runes.Equal([]rune("test"), []rune("TEST")), false)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("Test")), false)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("tEsT")), false)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("teSt")), false)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("test")), true)
	assert.Equal(t, runes.Equal([]rune("test"), []rune("test ")), false)
}

func TestIndex(t *testing.T) {
	assert.Equal(t, runes.Index([]rune("test"), []rune("test")), 0)
	assert.Equal(t, runes.Index([]rune("test"), []rune("test ")), -1)
	assert.Equal(t, runes.Index([]rune("test"), []rune("TEST")), -1)
	assert.Equal(t, runes.Index([]rune("test"), []rune("Test")), -1)
	assert.Equal(t, runes.Index([]rune("test"), []rune("tEsT")), -1)
	assert.Equal(t, runes.Index([]rune("test"), []rune("teSt")), -1)
}

func TestIndexFold(t *testing.T) {
	assert.Equal(t, 0, runes.IndexFold([]rune("test"), []rune("test")))
	assert.Equal(t, 0, runes.IndexFold([]rune("test"), []rune("TEST")))
	assert.Equal(t, 0, runes.IndexFold([]rune("test"), []rune("Test")))
	assert.Equal(t, 0, runes.IndexFold([]rune("test"), []rune("tEsT")))
	assert.Equal(t, 0, runes.IndexFold([]rune("test"), []rune("teSt")))
	assert.Equal(t, -1, runes.IndexFold([]rune("test"), []rune("test ")))
	assert.Equal(t, 1, runes.IndexFold([]rune(" test "), []rune("teSt")))
}
