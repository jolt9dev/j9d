//go:build aix || darwin || dragonfly || freebsd || hurd || illumos || ios || linux || netbsd || openbsd || plan9 || solaris || zos
// +build aix darwin dragonfly freebsd hurd illumos ios linux netbsd openbsd plan9 solaris zos

package env

import (
	"fmt"

	"github.com/jolt9dev/j9d/pkg/platform"
)

const (
	// The path variable name for the current OS.
	PATH = "PATH"
	// The home directory variable name for the current OS.
	HOME = "HOME"
	// The host name variable name for the current OS.
	HOSTNAME = "HOSTNAME"
	// The user name variable name for the current OS.
	USER = "USER"
	// The temporary directory for the current user. The variable
	// may not be defined on all systems.
	TMP = "TMPDIR"
	// The home config directory for the current user. The variable
	// may not be defined on all systems.
	HOME_CONFIG = "XDG_CONFIG_HOME"
	// The home data directory for the current user. The variable
	// may not be defined on all systems.
	HOME_DATA = "XDG_DATA_HOME"
	// The home cache directory for the current user. The variable
	// may not be defined on all systems.
	HOME_CACHE = "XDG_CACHE_HOME"
)

func Getx(key string, x int) string {
	switch x {
	case X_PROCESS:
		return Get(key)
	case X_MACHINE:
		return Get(key)
	case X_USER:
		return Get(key)

	default:
		return ""
	}
}

func Setx(key, value string, x int) error {
	switch x {
	case X_PROCESS:
		return Set(key, value)
	case X_MACHINE:
		return platform.ErrOsNotSupported
	case X_USER:
		return platform.ErrOsNotSupported
	}

	return fmt.Errorf("unknown x value: %d", x)
}

func Deletex(key string, x int) error {
	switch x {
	case X_PROCESS:
		return Delete(key)
	case X_MACHINE:
		return platform.ErrOsNotSupported
	case X_USER:
		return platform.ErrOsNotSupported
	}

	return fmt.Errorf("unknown x value: %d", x)
}

func hasPath(path string, paths []string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}
