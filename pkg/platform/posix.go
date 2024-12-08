//go:build aix || darwin || dragonfly || freebsd || hurd || illumos || ios || linux || netbsd || openbsd || plan9 || solaris || zos

package platform

import (
	o "os"
)

const (
	POSIX    = true
	EOL      = "\n"
	PATH_SEP = ":"
	DIR_SEP  = "/"
	DEV_NULL = "/dev/null"
)

func IsWsl() bool {
	_, err := o.Stat("/proc/sys/fs/binfmt_misc/WSLInterop")
	return err == nil
}

func IsWindows() bool {
	return false
}
