//go:build windows

package ospaths

import (
	"path/filepath"

	"github.com/jolt9dev/j9d/pkg/env"
	"golang.org/x/sys/windows"
)

func HomeDir() (string, error) {
	home := env.Get(env.HOME)
	if home == "" {
		return windows.KnownFolderPath(windows.FOLDERID_Profile, windows.KF_FLAG_DEFAULT)
	}

	return home, nil
}

func HomeCacheDir() (string, error) {
	localAppData := env.Get(env.HOME_CACHE)
	if localAppData != "" {
		return localAppData, nil
	}

	localAppData, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		localAppData = filepath.Join(homeDir, "AppData", "Local")
		return localAppData, nil
	}

	return localAppData, nil
}

func HomeConfigDir() (string, error) {
	roamingAppData := env.Get(env.HOME_CONFIG)
	if roamingAppData != "" {
		return roamingAppData, nil
	}

	roamingAppData, err := windows.KnownFolderPath(windows.FOLDERID_RoamingAppData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		roamingAppData = filepath.Join(homeDir, "AppData", "Roaming")
		return roamingAppData, nil
	}

	return roamingAppData, nil
}

func HomeDataDir() (string, error) {
	localAppData := env.Get(env.HOME_DATA)
	if localAppData != "" {
		return localAppData, nil
	}

	localAppData, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}

		localAppData = filepath.Join(homeDir, "AppData", "Local")
		return localAppData, nil
	}

	return localAppData, nil
}

func HomeBinDir() (string, error) {
	localAppData, err := HomeDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(localAppData, "Programs", "bin"), nil
}

func AppHomeConfigDir(appName string) (string, error) {
	configDir, err := HomeConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, appName), nil
}

func AppHomeDataDir(appName string) (string, error) {
	dataDir, err := HomeDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, appName, "data"), nil
}

func AppHomeCacheDir(appName string) (string, error) {
	cacheDir, err := HomeCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, appName, "cache"), nil
}

func OsConfigDir() (string, error) {
	path := env.Get("OS_CONFIG_DIR")
	if path != "" {
		return path, nil
	}

	path = env.Get("ProgramData")
	if path != "" {
		env.Set("OS_CONFIG_DIR", path)
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	path = dir
	env.Set("OS_CONFIG_DIR", path)

	return path, nil
}

func OsDataDir() (string, error) {
	path := env.Get("OS_DATA_DIR")
	if path != "" {
		return path, nil
	}

	path = env.Get("ProgramData")
	if path != "" {
		env.Set("OS_DATA_DIR", path)
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	path = dir
	env.Set("OS_DATA_DIR", path)

	return path, nil
}

func OsCacheDir() (string, error) {
	path := env.Get("OS_CACHE_DIR")
	if path != "" {
		return path, nil
	}

	path = env.Get("ProgramData")
	if path != "" {
		path = filepath.Join(path, "cache")
		env.Set("OS_CACHE_DIR", path)
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	path = filepath.Join(dir, "cache")
	env.Set("OS_DATA_DIR", path)

	return path, nil
}

func OsBinDir() (string, error) {
	path := env.Get("OS_BIN_DIR")
	if path != "" {
		return path, nil
	}

	path = env.Get("ProgramFiles")
	if path != "" {
		path = filepath.Join(path, "bin")
		env.Set("OS_BIN_DIR", path)
		return path, nil
	}

	path = env.Get("ProgramFiles(x86)")
	if path != "" {
		path = filepath.Join(path, "bin")
		env.Set("OS_BIN_DIR", path)
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	path = filepath.Join(dir, "bin")
	env.Set("OS_BIN_DIR", path)

	return path, nil
}

func programData() (string, error) {
	path := env.Get("PROGRAMDATA")
	if path != "" {
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func AppCacheDir(appName string) (string, error) {
	pd, err := programData()
	if err != nil {
		return "", err
	}

	return filepath.Join(pd, appName, "cache"), nil
}

func AppConfigDir(appName string) (string, error) {
	pd, err := programData()
	if err != nil {
		return "", err
	}

	return filepath.Join(pd, appName, "config"), nil
}

func AppDataDir(appName string) (string, error) {
	pd, err := programData()
	if err != nil {
		return "", err
	}

	return filepath.Join(pd, appName, "data"), nil
}

func OptDir() (string, error) {
	path := env.Get("OPT_DIR")
	if path != "" {
		return path, nil
	}

	path = env.Get("ProgramFiles")
	if path != "" {
		env.Set("OPT_DIR", path)
		return path, nil
	}

	path = env.Get("ProgramFiles(x86)")
	if path != "" {
		env.Set("OPT_DIR", path)
		return path, nil
	}

	dir, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	path = dir
	env.Set("OPT_DIR", path)

	return path, nil
}

func HomeDocumentsDir() (string, error) {
	dir, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func HomeDownloadsDir() (string, error) {
	dir, err := windows.KnownFolderPath(windows.FOLDERID_Downloads, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	return dir, nil
}
