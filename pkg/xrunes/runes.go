package runes

import (
	"slices"
	"unicode"
)

func Contains(s []rune, r []rune) bool {
	return Index(s, r) > -1
}

func ContainsFold(s []rune, r []rune) bool {
	return IndexFold(s, r) > -1
}

func Equal(x []rune, y []rune) bool {
	return slices.Equal(x, y)
}

func EqualFold(x []rune, y []rune) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			if unicode.IsLetter(x[i]) {
				if equalFoldRune(x[i], y[i]) {
					continue
				}
			}

			return false
		}
	}

	return true
}

func HasPrefix(s []rune, prefix []rune) bool {
	if len(s) < len(prefix) {
		return false
	}

	return slices.Equal(s[:len(prefix)], prefix)
}

func HasPrefixFold(s []rune, prefix []rune) bool {
	sl := len(s)
	rl := len(prefix)

	if rl == 0 {
		return true
	}

	if sl < rl {
		return false
	}

	return EqualFold(s[:rl], prefix)
}

func HasSuffix(s []rune, suffix []rune) bool {
	if len(s) < len(suffix) {
		return false
	}

	return slices.Equal(s[len(s)-len(suffix):], suffix)
}

func HasSuffixFold(s []rune, suffix []rune) bool {
	sl := len(s)
	rl := len(suffix)

	if rl == 0 {
		return true
	}

	if sl < rl {
		return false
	}

	return EqualFold(s[sl-rl:], suffix)
}

func IndexRune(s []rune, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}

	// Not found
	return -1
}

func IndexRuneFold(s []rune, r rune) int {
	if len(s) == 0 {
		return -1
	}

	for i, c := range s {
		if c == r {
			return i
		}

		if unicode.IsLetter(c) {
			if equalFoldRune(c, r) {
				return i
			}
		}
	}

	return -1
}

func Index(s []rune, r []rune) int {
	sl := len(s)
	rl := len(r)
	if rl == 0 {
		return 0
	}

	if sl < rl {
		return -1
	}

	for i := 0; i < sl; i++ {
		if i+rl > sl {
			return -1
		}

		for j, y := range r {
			x := s[i+j]
			if x == y {
				if j == rl-1 {
					return i
				}

				continue
			}

			break
		}
	}

	return -1
}

func IndexFold(s []rune, r []rune) int {
	sl := len(s)
	rl := len(r)
	if rl == 0 {
		return 0
	}

	for i := 0; i < sl; i++ {
		if i+(rl) > sl {
			return -1
		}

		for j, y := range r {
			x := s[i+j]
			if x == y {
				if j == rl-1 {
					return i
				}

				continue
			}

			if unicode.IsLetter(x) {
				if equalFoldRune(x, y) {
					if j == rl-1 {
						return i
					}

					continue
				}
			}

			break
		}
	}

	return -1
}

func Trim(s []rune, cutset []rune) []rune {
	if len(s) == 0 {
		return s
	}

	start := 0
	end := len(s)

	for i := 0; i < len(s); i++ {
		if slices.Contains(cutset, s[i]) {
			start = i + 1
			break
		}
	}

	i := len(s) - 1
	for i >= 0 {
		if slices.Contains(cutset, s[i]) {
			end = i
			break
		}

		i--
	}

	return s[start:end]
}

func TrimLeft(s []rune, cutset []rune) []rune {
	if len(s) == 0 {
		return s
	}

	if len(cutset) == 0 {
		return s
	}

	start := 0

	for i := 0; i < len(s); i++ {
		if slices.Contains(cutset, s[i]) {
			start = i + 1
			break
		}
	}

	// if start == len(s) {
	// 	return s
	// }

	return s[start:]
}

func TrimRight(s []rune, cutset []rune) []rune {
	if len(s) == 0 {
		return s
	}

	if len(cutset) == 0 {
		return s
	}

	end := len(s)

	for i := len(s) - 1; i >= 0; i-- {
		if !slices.Contains(cutset, s[i]) {
			end = i
			break
		}
	}

	// if end == 0 {
	// 	return s
	// }

	return s[:end]
}

func equalFoldRune(x, y rune) bool {
	xx := unicode.SimpleFold(x)
	if xx == y {
		return true
	}
	yy := unicode.SimpleFold(y)
	return yy == x
}
