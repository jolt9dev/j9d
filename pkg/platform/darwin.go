//go:build darwin || ios || tvos || watchos

package platform

const (
	FAMILY = "darwin"
)

func IsDarwin() bool {
	return true
}
