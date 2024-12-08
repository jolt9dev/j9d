package xstrings

import (
	str "strings"
)

func PadLeft(s string, length int, pad string) string {
	diff := length - len(s)
	if diff <= 0 {
		return s
	}

	return str.Repeat(pad, diff) + s
}

func PadRight(s string, length int, pad string) string {
	diff := length - len(s)
	if diff <= 0 {
		return s
	}

	return s + str.Repeat(pad, diff)
}
