package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptValueCompare(t *testing.T) {
	out, gotError := ParseBcryptValue("10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue.")
	assert.NoError(t, gotError)

	got := BcryptValueCompare("password", out)
	assert.True(t, got)
}
