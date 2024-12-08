package strings_test

import (
	"testing"

	strings "github.com/jolt9dev/j9d/pkg/xstrings"
	"github.com/stretchr/testify/assert"
)

func TestEmptySpace(t *testing.T) {
	assert.True(t, strings.IsEmptySpace(""))
	assert.True(t, strings.IsEmptySpace(" "))
	assert.True(t, strings.IsEmptySpace("  "))
	assert.True(t, strings.IsEmptySpace(" \n\t\v\f\r"))
}
