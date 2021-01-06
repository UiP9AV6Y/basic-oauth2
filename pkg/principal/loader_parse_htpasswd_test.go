package principal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/password"
)

func TestParseHtpasswd(t *testing.T) {
	out := NewPassthroughLoader()
	testCases := []parseFileTestCase{
		parseFileTestCase{
			err:    "testdata/parse_htpasswd/empty_ident.htpasswd (1): principal ident cannot be empty",
			file:   "testdata/parse_htpasswd/empty_ident.htpasswd",
			result: nil,
		},
		parseFileTestCase{
			file:   "testdata/parse_htpasswd/empty.htpasswd",
			result: []Principal{},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/string.htpasswd",
			result: []Principal{
				&DefaultPrincipal{
					Ident:    "string",
					Email:    newTestEmail("string"),
					Password: password.NewPlaintextPassword("password"),
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/with_email.htpasswd",
			result: []Principal{
				&DefaultPrincipal{
					Ident:    "with_email",
					Email:    "user@example.org",
					Password: password.NewPlaintextPassword("password"),
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/email_ident.htpasswd",
			result: []Principal{
				&DefaultPrincipal{
					Ident:         "email@example.net",
					Email:         "email@example.net",
					Password:      password.NewPlaintextPassword("password"),
					EmailVerified: true,
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/email_ident_email.htpasswd",
			result: []Principal{
				&DefaultPrincipal{
					Ident:         "email@example.com",
					Password:      password.NewPlaintextPassword("password"),
					Email:         "user@example.dev",
					EmailVerified: true,
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/bcrypt.htpasswd",
			result: []Principal{
				newTestPrincipal("bcrypt",
					"$2y$10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue."),
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/sha1.htpasswd",
			result: []Principal{
				newTestPrincipal("sha1",
					"{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g="),
			},
		},
		parseFileTestCase{
			file: "testdata/parse_htpasswd/ident.htpasswd",
			result: []Principal{
				newTestPrincipal("ident", "ident"),
			},
		},
	}

	for _, c := range testCases {
		runner := func(t *testing.T) {
			gotPrincipals, gotError := out.ParseHtpasswd(c.file)

			if c.err == "" {
				assert.NoError(t, gotError)
			} else {
				assert.EqualError(t, gotError, c.err)
			}

			if c.result == nil {
				assert.Nil(t, gotPrincipals)
			} else if assert.NotNil(t, gotPrincipals) {
				comparePrincipals(t, c.result, gotPrincipals)
			}
		}

		t.Run(c.file, runner)
	}
}
