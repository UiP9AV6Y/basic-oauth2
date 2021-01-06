package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPlaintextValue(t *testing.T) {
	testCases := map[string]bool{
		"test":   true,
		"123":    true,
		"TEST":   true,
		"!test":  false,
		"$test":  false,
		"_test":  false,
		"":       false,
		"{test}": false,
	}

	for input, want := range testCases {
		runner := func(t *testing.T) {
			got := IsPlaintextValue(input)

			assert.Equal(t, want, got, input)
		}

		t.Run(input, runner)
	}
}
