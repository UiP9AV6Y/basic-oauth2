package principal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/password"
)

type parseFileTestCase struct {
	err    string
	file   string
	result []Principal
}

func TestParseFile(t *testing.T) {
	out := NewPassthroughLoader()
	testCases := []parseFileTestCase{
		parseFileTestCase{
			err:    "principal ident cannot be empty",
			file:   "testdata/parse_file/empty_ident.yaml",
			result: nil,
		},
		parseFileTestCase{
			file: "testdata/parse_file/string.yaml",
			result: []Principal{
				&DefaultPrincipal{
					Ident:    "string",
					Email:    newTestEmail("string"),
					Password: password.NewPlaintextPassword("password"),
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/with_email.yaml",
			result: []Principal{
				&DefaultPrincipal{
					Ident:    "with_email",
					Email:    "user@example.org",
					Password: password.NewPlaintextPassword("password"),
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/email_ident.yaml",
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
			file: "testdata/parse_file/email_ident_email.yaml",
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
			file: "testdata/parse_file/empty_mapping.yaml",
			result: []Principal{
				&DefaultPrincipal{
					Ident:    "empty_mapping",
					Email:    newTestEmail("empty_mapping"),
					Password: password.NewPlaintextPassword("empty_mapping"),
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/bcrypt.yaml",
			result: []Principal{
				newTestPrincipal("bcrypt",
					"$2y$10$Hr.ji.jd4yKMr3nscluW/.3VFCtLWD/zPUjrpTbB/fFbUv.z6Sue."),
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/sha1.yaml",
			result: []Principal{
				newTestPrincipal("sha1",
					"{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g="),
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/generated_verified.yaml",
			result: []Principal{
				&DefaultPrincipal{
					Ident:         "generated_verified",
					Password:      password.NewPlaintextPassword("generated_verified"),
					Email:         newTestEmail("generated_verified"),
					EmailVerified: true,
				},
			},
		},
		parseFileTestCase{
			file: "testdata/parse_file/email_verified.yaml",
			result: []Principal{
				&DefaultPrincipal{
					Ident:         "email_verified",
					Password:      password.NewPlaintextPassword("email_verified"),
					Email:         "user@example.io",
					EmailVerified: true,
				},
			},
		},
	}

	for _, c := range testCases {
		runner := func(t *testing.T) {
			gotPrincipals, gotError := out.ParseFile(c.file)

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

func newTestEmail(ident string) (email string) {
	email, _ = GenerateEmail(ident, DefaultEmailDomain)
	return
}

func newTestPrincipal(ident, pass string) Principal {
	principal := &DefaultPrincipal{
		Ident: ident,
	}
	principal.Password, _ = password.ParsePassword(pass)
	principal.Email, principal.EmailVerified = GenerateEmail(ident, DefaultEmailDomain)

	return principal
}

func comparePrincipals(t *testing.T, want, got []Principal) {
	assert.Len(t, got, len(want), "Principals")

	for _, w := range want {
		var found Principal

		for _, g := range got {
			if w.GetIdent() == g.GetIdent() {
				found = g
				break
			}
		}

		if assert.NotNil(t, found, "no principal found with Ident=%q", w.GetIdent()) {
			comparePrincipal(t, w, found)
		}
	}
}

func comparePrincipal(t *testing.T, want, got Principal) {
	wantEmail, wantVerified := want.GetEmail()
	gotEmail, gotVerified := got.GetEmail()

	assert.Equal(t, want.GetIdent(), got.GetIdent(), "Ident")
	assert.Equal(t, wantEmail, gotEmail, "Email")
	assert.Equal(t, wantVerified, gotVerified, "Verified")
}
