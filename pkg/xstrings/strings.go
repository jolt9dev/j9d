package xstrings

import (
	str "strings"
)

const TEST = "TEST"

func init() {
}

func ToLower(s string) string {
	return str.ToLower(s)
}

func ToUpper(s string) string {
	return str.ToUpper(s)
}

func TrimSpace(s string) string {
	return str.TrimSpace(s)
}

func Split(s, sep string) []string {
	return str.Split(s, sep)
}

func SplitAfter(s, sep string) []string {
	return str.SplitAfter(s, sep)
}

func SplitAfterN(s, sep string, n int) []string {
	return str.SplitAfterN(s, sep, n)
}

func SplitN(s, sep string, n int) []string {
	return str.SplitN(s, sep, n)
}

func Join(elems []string, sep string) string {
	return str.Join(elems, sep)
}

func Contains(s, substr string) bool {
	return str.Contains(s, substr)
}

func ContainsAny(s, chars string) bool {
	return str.ContainsAny(s, chars)
}

func ContainsRune(s string, r rune) bool {
	return str.ContainsRune(s, r)
}

func ContainsFunc(s string, f func(rune) bool) bool {
	return str.ContainsFunc(s, f)
}

func Count(s, substr string) int {
	return str.Count(s, substr)
}

func EqualFold(s, t string) bool {
	return str.EqualFold(s, t)
}

func Fields(s string) []string {
	return str.Fields(s)
}

func FieldsFunc(s string, f func(rune) bool) []string {
	return str.FieldsFunc(s, f)
}

func Index(s, substr string) int {
	return str.Index(s, substr)
}

func IndexAny(s, chars string) int {
	return str.IndexAny(s, chars)
}

func IndexByte(s string, c byte) int {
	return str.IndexByte(s, c)
}

func IndexFunc(s string, f func(rune) bool) int {
	return str.IndexFunc(s, f)
}

func IndexRune(s string, r rune) int {

	return str.IndexRune(s, r)
}

func SplitAny(s, sep string) []string {
	set := make([]string, 0)
	sb := str.Builder{}
	for _, r := range s {
		if str.ContainsRune(sep, r) {
			set = append(set, sb.String())
			sb.Reset()
			continue
		}

		sb.WriteRune(r)
	}

	if sb.Len() > 0 {
		set = append(set, sb.String())
	}

	return set
}
