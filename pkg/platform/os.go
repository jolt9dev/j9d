package platform

import (
	"fmt"
	"runtime"
)

const (
	OS   = runtime.GOOS
	ARCH = runtime.GOARCH
)

var (
	ErrOsNotSupported = fmt.Errorf("os %s not supported", runtime.GOOS)
)

func init() {
}
