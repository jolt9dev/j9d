//go:build windows

package cps

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type elevationData struct {
	TokenIsElevated int32
}

func queryTokenElevation() (bool, error) {

	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)

	if err != nil {
		return false, err
	}

	defer token.Close()

	info := make([]byte, 4)
	var returnLength uint32

	err = windows.GetTokenInformation(token, windows.TokenElevation, &info[0], uint32(4), &returnLength)
	if err != nil {
		return false, err
	}

	elevation := (*elevationData)(unsafe.Pointer(&info[0]))
	return elevation.TokenIsElevated != 0, nil
}

func IsElevated() bool {
	isElevated, err := queryTokenElevation()
	if err != nil {
		return false
	}

	return isElevated
}
