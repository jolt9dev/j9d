package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/jolt9dev/j9d/pkg/platform"
)

const (
	X_PROCESS = 0
	X_MACHINE = 1
	X_USER    = 2
)

func Get(key string) string {
	return os.Getenv(key)
}

func Set(key, value string) error {
	return os.Setenv(key, value)
}

func Delete(key string) error {
	return os.Unsetenv(key)
}

func Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func All() map[string]string {
	kv := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) == 2 && len(pair[1]) > 0 {
			kv[pair[0]] = pair[1]
		}
	}

	return kv
}

func Print() {
	for k, v := range All() {
		fmt.Printf("%s=%s%s", k, v, platform.EOL)
	}
}

func GetPath() string {
	return os.Getenv(PATH)
}

func SetPath(path string) error {
	return os.Setenv(PATH, path)
}

func HasPath(path string) bool {
	return hasPath(path, SplitPath())
}

func AppendPath(path string) error {
	paths := SplitPath()
	if hasPath(path, paths) {
		return nil
	}
	paths = append(paths, path)
	return SetPath(JoinPath(paths...))
}

func PrependPath(path string) error {
	paths := SplitPath()
	if hasPath(path, paths) {
		return nil
	}
	paths = append([]string{path}, paths...)
	return SetPath(JoinPath(paths...))
}

func SplitPath() []string {
	return strings.Split(GetPath(), string(os.PathListSeparator))
}

func JoinPath(paths ...string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}
