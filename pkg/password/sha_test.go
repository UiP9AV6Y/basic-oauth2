package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHAGenerate(t *testing.T) {
	out, gotError := ParseSHAValue("W6ph5Mm5Pz8GgiULbPgzG37mj9g=")
	assert.NoError(t, gotError)

	want := "{SHA}JJujYAACm76XSZwD21qQAfa3NOw="
	got, gotError2 := out.Generate("username")
	assert.NoError(t, gotError2)
	assert.Equal(t, want, got)
}

func TestSHAString(t *testing.T) {
	out, gotError := ParseSHAValue("W6ph5Mm5Pz8GgiULbPgzG37mj9g=")
	assert.NoError(t, gotError)

	want := "{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g="
	got := out.String()
	assert.Equal(t, want, got)
}
