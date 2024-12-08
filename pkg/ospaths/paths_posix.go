//go:build aix || darwin || dragonfly || freebsd || hurd || illumos || ios || linux || netbsd || openbsd || plan9 || solaris || zos

package ospaths

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/jolt9dev/j9d/pkg/env"
)

func OptDir() (string, error) {
	if runtime.GOOS == "darwin" {
		return "/Applications", nil
	}

	return "/opt", nil
}

func HomeDir() (string, error) {
	homeDir := env.Get(env.HOME)
	if homeDir != "" {
		return homeDir, nil
	}

	homeDir = env.Get("XDG_CONFIG_HOME")
	if homeDir != "" {
		env.Set(env.HOME, homeDir)
		return homeDir, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir = usr.HomeDir
	env.Set(env.HOME, homeDir)
	return homeDir, nil
}

func HomeBinDir() (string, error) {
	binDir := env.Get("XDG_BIN_HOME")
	if binDir != "" {
		return binDir, nil
	}

	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}

	binDir = filepath.Join(homeDir, ".local", "bin")
	env.Set("XDG_BIN_HOME", binDir)
	return binDir, nil
}

func HomeCacheDir() (string, error) {
	cacheDir := env.Get(env.HOME_CACHE)
	if cacheDir != "" {
		return cacheDir, nil
	}

	if runtime.GOOS == "darwin" {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		cacheDir = filepath.Join(homeDir, "Library", "Caches")
		env.Set(env.HOME_CACHE, cacheDir)
		return cacheDir, nil
	}

	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}

	cacheDir = filepath.Join(homeDir, ".cache")
	env.Set(env.HOME_CACHE, cacheDir)
	return cacheDir, nil
}

func HomeConfigDir() (string, error) {
	configDir := env.Get(env.HOME_CONFIG)
	if configDir != "" {
		return configDir, nil
	}

	if runtime.GOOS == "darwin" {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		configDir = filepath.Join(homeDir, "Library", "Application Support")
		env.Set(env.HOME_CONFIG, configDir)
		return configDir, nil
	}

	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}

	os.UserConfigDir()
	configDir = filepath.Join(homeDir, ".config")
	env.Set(env.HOME_CONFIG, configDir)
	return configDir, nil
}

func HomeDataDir() (string, error) {
	dataDir := env.Get(env.HOME_DATA)
	if dataDir != "" {
		return dataDir, nil
	}

	if runtime.GOOS == "darwin" {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		dataDir = filepath.Join(homeDir, "Library", "Application Support")
		env.Set(env.HOME_DATA, dataDir)
		return dataDir, nil
	}

	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}

	dataDir = filepath.Join(homeDir, ".local", "share")
	env.Set(env.HOME_DATA, dataDir)
	return dataDir, nil
}

func AppHomeConfigDir(appName string) (string, error) {
	configDir, err := HomeConfigDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "darwin" {
		configDir = filepath.Join(configDir, appName, "config")
		return configDir, nil
	}

	return filepath.Join(configDir, appName), nil
}

func AppHomeDataDir(appName string) (string, error) {
	dataDir, err := HomeDataDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "darwin" {
		dataDir = filepath.Join(dataDir, appName, "data")
		return dataDir, nil
	}

	return filepath.Join(dataDir, appName), nil
}

func AppHomeCacheDir(appName string) (string, error) {
	cacheDir, err := HomeCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, appName), nil
}

func AppConfigDir(appName string) (string, error) {
	return filepath.Join("/etc", appName), nil
}

func AppDataDir(appName string) (string, error) {
	return filepath.Join("/usr/local/share", appName), nil
}

func AppCacheDir(appName string) (string, error) {
	return filepath.Join("/var/cache", appName), nil
}

func OsBinDir(appName string) (string, error) {
	return "/usr/local/bin", nil
}

func HomeDocumentsDir() (string, error) {
	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "Documents"), nil
}

func HomeDownloadsDir() (string, error) {
	homeDir, err := HomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "Downloads"), nil
}
