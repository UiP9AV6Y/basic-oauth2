package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePassword(t *testing.T) {
	// "password" (encoded) -> codec name
	testCases := map[string]string{
		"password":                          "Plaintext",
		"{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g=": "SHA",
		"$2y$10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue.": "Bcrypt",
	}

	for input, wantCodecName := range testCases {
		runner := func(t *testing.T) {
			gotPassword, gotError := ParsePassword(input)

			assert.NoError(t, gotError)

			if assert.NotNil(t, gotPassword) {
				assert.Equal(t, wantCodecName, gotPassword.Codec.Name(), "Codec")
			}
		}

		t.Run(input, runner)
	}
}

func TestPasswordCompare(t *testing.T) {
	plainV, plainErr := ParsePassword("password")
	shaV, shaErr := ParsePassword("{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g=")
	bcryptV, bcryptErr := ParsePassword("$2y$10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue.")

	assert.NoError(t, plainErr, "Plain value parsing")
	assert.NoError(t, shaErr, "SHA value parsing")
	assert.NoError(t, bcryptErr, "Bcrypt value parsing")

	assert.True(t, plainV.Compare("password"), "Plain")
	assert.True(t, shaV.Compare("password"), "SHA")
	assert.True(t, bcryptV.Compare("password"), "Bcrypt")
}
