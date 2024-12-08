package env_test

import (
	"testing"

	"github.com/jolt9dev/j9d/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	home := ""
	home = env.Get(env.HOME)

	assert.NotEmpty(t, home)
}

func TestSet(t *testing.T) {
	env.Set("SET_TEST", "test")
	test := env.Get("SET_TEST")

	if len(test) == 0 {
		t.Errorf("Expected %v, got %v", "non empty sring for test", test)
	}

	if test != "test" {
		t.Errorf("Expected %v, got %v", "test", test)
	}
}

func TestDelete(t *testing.T) {
	env.Set("SET_TEST2", "test")
	test := env.Get("SET_TEST2")

	if test != "test" {
		t.Errorf("Expected %v, got %v", "test", test)
	}

	env.Delete("SET_TEST2")
	test = env.Get("SET_TEST2")
	if test != "" {
		t.Errorf("Expected empty string, got %v", test)
	}
}

func TestPaths(t *testing.T) {
	envPath := env.GetPath()

	if envPath == "" {
		t.Errorf("Expected non empty value, got %v", envPath)
	}
}

func TestSplitPaths(t *testing.T) {
	envPaths := env.SplitPath()

	if len(envPaths) == 0 {
		t.Errorf("Expected more than zero length, got %v", len(envPaths))
	}
}
