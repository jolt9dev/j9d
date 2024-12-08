package xexec_test

import (
	"strings"
	"testing"

	"github.com/jolt9dev/j9d/pkg/xexec"
	"github.com/stretchr/testify/assert"
)

func TestNewCommandOutput(t *testing.T) {
	echo, ok := xexec.Which("echo")
	if !ok {
		t.Skip("echo not found")
	}

	o, err := xexec.New(echo, "hello").Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "hello", strings.TrimSpace(o.Text()))
}

func TestCommandOutput(t *testing.T) {
	_, ok := xexec.Which("echo")
	if !ok {
		t.Skip("echo not found")
	}

	cmd := "echo 'hello world'"

	o, err := xexec.Command(cmd).Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "hello world", strings.TrimSpace(o.Text()))
}

func TestPipeCommand(t *testing.T) {
	_, hasGrep := xexec.Which("grep")
	_, hasEcho := xexec.Which("echo")

	if !hasEcho || !hasGrep {
		t.Skip("grep or echo not found")
	}

	o, err := xexec.Command("echo 'Hello World'").PipeCommand("grep Hello").Output()
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Code)
	assert.Equal(t, "Hello World", strings.TrimSpace(o.Text()))
}
