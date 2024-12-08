package platform_test

import (
	"runtime"
	"testing"

	"github.com/jolt9dev/j9d/pkg/platform"
	"github.com/stretchr/testify/assert"
)

func TestOs(t *testing.T) {
	assert.Equal(t, platform.OS, runtime.GOOS)
}

func TestWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		assert.True(t, platform.IsWindows())
		return
	}

	assert.False(t, platform.IsWindows())
}
