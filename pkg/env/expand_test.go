package env_test

import (
	"os"
	"testing"

	"github.com/jolt9dev/j9d/pkg/env"
)

func TestExpandNoReplace(t *testing.T) {
	out1, err := env.Expand("test", &env.ExpandOptions{})

	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if out1 != "test" {
		t.Errorf("Expected %s, got %s", "test", out1)
	}
}

func TestExpand(t *testing.T) {
	os.Setenv("WORLD", "World")

	out1, err := env.Expand("Hello $WORLD", &env.ExpandOptions{})

	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if out1 != "Hello World" {
		t.Errorf("Expected %s, got %s", "Hello World", out1)
	}
}

func TestExpandWithParen(t *testing.T) {
	os.Setenv("WORLD", "World")

	out1, err := env.Expand("Hello ${WORLD}", &env.ExpandOptions{})

	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if out1 != "Hello World" {
		t.Errorf("Expected %s, got %s", "Hello World", out1)
	}
}

func TestExpandWithDefault(t *testing.T) {
	os.Setenv("WORLD", "Emma")

	out1, err := env.Expand("Hello ${Bad:-World}", &env.ExpandOptions{})

	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if out1 != "Hello World" {
		t.Errorf("Expected %s, got %s", "Hello World", out1)
	}
}

func TestExpandWithErrorMessage(t *testing.T) {
	os.Setenv("WORLD", "Emma")

	out1, err := env.Expand("Hello ${Bad:?Error}", &env.ExpandOptions{})
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}

	if out1 != "" {
		t.Errorf("Expected %s, got %s", "", out1)
	}
}
