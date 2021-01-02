package utils

import (
	"strings"
)

func MapFields(s string) map[string]bool {
	f := strings.Fields(s)
	m := make(map[string]bool, len(f))

	for _, s := range f {
		m[s] = true
	}

	return m
}
