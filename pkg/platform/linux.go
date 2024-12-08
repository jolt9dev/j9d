//go:build linux || android || solaris || illumos || plan9

package platform

const (
	FAMILY = "linux"
)

func IsDarwin() bool {
	return false
}
