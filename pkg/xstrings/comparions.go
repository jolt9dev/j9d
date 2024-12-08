package xstrings

import (
	str "strings"
	"unicode"

	runes "github.com/jolt9dev/j9d/pkg/xrunes"
)

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsEmptySpace(s string) bool {
	if len(s) == 0 {
		return true
	}

	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

func HasPrefix(s, prefix string) bool {
	return str.HasPrefix(s, prefix)
}

func HasPrefixFold(s, prefix string) bool {
	return runes.HasPrefixFold([]rune(s), []rune(prefix))
}

func HasSuffix(s, suffix string) bool {
	return str.HasSuffix(s, suffix)
}

func HasSuffixFold(s, suffix string) bool {
	return runes.HasSuffixFold([]rune(s), []rune(suffix))
}
