package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCodec(t *testing.T) {
	// "password" (encoded) -> codec name
	testCases := map[string]string{
		"password":                          "Plaintext",
		"{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g=": "SHA",
		"$2y$10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue.": "Bcrypt",
	}

	for input, wantCodecName := range testCases {
		runner := func(t *testing.T) {
			gotCodec, gotParams, gotError := FindCodec(input)
			assert.NoError(t, gotError)
			assert.NotEmpty(t, gotParams, "Params")

			if assert.NotNil(t, gotCodec) {
				assert.Equal(t, wantCodecName, gotCodec.Name(), "Codec")
			}
		}

		t.Run(input, runner)
	}
}
