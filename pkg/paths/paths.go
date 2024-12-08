package paths

import (
	"path/filepath"

	p "github.com/jolt9dev/j9d/pkg/ospaths"
)

func ConfigDir() (string, error) {
	dir, err := p.AppConfigDir("deploy")
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "deploy"), nil
}

func CacheDir() (string, error) {
	dir, err := p.AppCacheDir("deploy")
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "deploy"), nil
}

func DataDir() (string, error) {
	dir, err := p.AppDataDir("deploy")
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "deploy"), nil
}
