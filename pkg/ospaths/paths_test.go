package ospaths_test

import (
	"testing"

	"github.com/jolt9dev/j9d/pkg/ospaths"
	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	home, err := ospaths.HomeDir()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	assert.DirExists(t, home)
}
