package notes

import (
	"unicode/utf8"
)

func FirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func Ellipsis(s string, threshold int, mark string) string {
	if utf8.RuneCountInString(s) <= threshold {
		return s
	}
	if len(mark) == 0 {
		return FirstN(s, threshold)
	}
	return FirstN(s, threshold) + mark
}
