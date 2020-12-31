package utils

import (
	"math"
	"strings"
)

type StringReader string

func (r StringReader) Read(p []byte) (n int, err error) {
	pool := r.Bloat(len(p))

	for i, b := range []byte(pool) {
		p[i] = b
	}

	return len(p), nil
}

func (r StringReader) Bloat(size int) string {
	pool := string(r)

	if size == 0 {
		return ""
	} else if len(pool) < size {
		reps := math.Ceil(float64(size / len(pool)))
		pool = strings.Repeat(pool, int(reps))
	}

	return pool[0 : size-1]
}
