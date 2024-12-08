package cps_test

import (
	"testing"

	"github.com/jolt9dev/j9d/pkg/cps"
	"github.com/jolt9dev/j9d/pkg/platform"
	"github.com/stretchr/testify/assert"
)

func TestUserAndGroupIds(t *testing.T) {

	if platform.IsWindows() {
		t.Skip("Skipping test on windows")
	}

	assert.Greater(t, cps.Uid(), -1)
	assert.Greater(t, cps.Gid(), -1)

	assert.Greater(t, cps.Egid(), -1)
	assert.Greater(t, cps.Euid(), -1)
}

func TestCwd(t *testing.T) {
	cwd, err := cps.Cwd()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	assert.NotEmpty(t, cwd)
}

func TestProcessId(t *testing.T) {
	pid := cps.Pid()
	assert.Greater(t, pid, 0)

	ppid := cps.Ppid()
	assert.Greater(t, ppid, 0)
}

func TestWrite(t *testing.T) {
	b, err := cps.WriteString("test.txt")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	assert.Greater(t, b, 0)
}
