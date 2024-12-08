//go:build windows
// +build windows

//nolint

package ssh

import (
	"golang.org/x/crypto/ssh"
)

// monWinCh does nothing for now on windows
//
//nolint:all
func monWinCh(session *ssh.Session, fd uintptr) {
	// do nothing

	return
}
