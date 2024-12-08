package env

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	// The path variable name for the current OS.
	PATH = "Path"
	// The home directory variable name for the current OS.
	HOME = "UserProfile"
	// The host name variable name for the current OS.
	HOSTNAME = "COMPUTERNAME"
	// The user name variable name for the current OS.
	USER = "USERNAME"
	// The temporary directory for the current user. The variable
	// may not be defined on all systems.
	TMP = "TEMP"
	// The home config directory for the current user. The variable
	// may not be defined on all systems.
	HOME_CONFIG = "AppData"
	// The home data directory for the current user. The variable
	// may not be defined on all systems.
	HOME_DATA = "LocalAppData"
	// The home cache directory for the current user. The variable
	// may not be defined on all systems.
	HOME_CACHE = "LocalAppData"
)

func Setx(key, value string, x int) error {
	switch x {
	case X_PROCESS:
		return Set(key, value)
	case X_MACHINE:

		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
		if err != nil {
			return err
		}

		defer k.Close()
	case X_USER:
		k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
		if err != nil {
			return err
		}

		defer k.Close()

		return k.SetStringValue(key, value)

	}

	return fmt.Errorf("unknown x value: %d", x)
}

func Getx(key string, x int) string {
	switch x {
	case X_PROCESS:
		return Get(key)
	case X_MACHINE:
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
		if err != nil {
			return ""
		}

		defer k.Close()

		v, _, err := k.GetStringValue(key)
		if err != nil {
			return ""
		}

		return v
	case X_USER:
		k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
		if err != nil {
			return ""
		}

		defer k.Close()

		v, _, err := k.GetStringValue(key)
		if err != nil {
			return ""
		}

		return v

	default:
		return ""
	}
}

func Deletex(key string, x int) error {
	switch x {
	case X_PROCESS:
		return Delete(key)
	case X_MACHINE:
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
		if err != nil {
			return err
		}

		defer k.Close()

		return k.DeleteValue(key)
	case X_USER:
		k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
		if err != nil {
			return err
		}

		defer k.Close()

		return k.DeleteValue(key)
	}

	return fmt.Errorf("unknown x value: %d", x)
}

func hasPath(path string, paths []string) bool {
	for _, p := range paths {
		if strings.EqualFold(p, path) {
			return true
		}
	}
	return false
}
