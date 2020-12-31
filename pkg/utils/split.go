package utils

import "strings"

func SplitPair(s, delim string) (string, string) {
	split := strings.Index(s, delim)

	if split < 0 {
		return s, ""
	}

	return s[0:split], s[split+1:]
}

func SplitPrefix(s, prefix string) (string, bool) {
	sub := len(prefix)

	if len(s) <= sub || s[0:sub] != prefix {
		return "", false
	}

	return s[sub:], true
}
