package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringReaderTestCase struct {
	byteCount int
	rendered  string
}

func TestStringReaderRead(t *testing.T) {
	out := StringReader("abcd")
	testCases := []stringReaderTestCase{
		stringReaderTestCase{
			byteCount: 0,
		},
		stringReaderTestCase{
			byteCount: 2,
			rendered:  "ab",
		},
		stringReaderTestCase{
			byteCount: 9,
			rendered:  "abcdabcda",
		},
	}

	for _, c := range testCases {
		spec := fmt.Sprintf("reading %d bytes", c.byteCount)
		test := func(st *testing.T) {
			buffer := make([]byte, c.byteCount)
			length, err := out.Read(buffer)

			assert.Nil(t, err, "read error")
			assert.Equal(t, c.byteCount, length, "read length")
			assert.Equal(t, c.rendered, string(buffer), "result")
		}

		t.Run(spec, test)
	}
}
