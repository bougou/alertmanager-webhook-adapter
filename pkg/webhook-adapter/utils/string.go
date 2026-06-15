package utils

import "strings"

func TruncateToValidUTF8(s string, byteSize int, mark string) string {
	b := []byte(s)
	if len(b) > byteSize {
		l := byteSize - len(mark)
		return strings.ToValidUTF8(string(b[:l]), "") + mark
	}
	return s
}
