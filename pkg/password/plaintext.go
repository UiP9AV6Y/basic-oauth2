package password

import (
	"unicode"
	"unicode/utf8"
)

var plaintextCodec = &CodecDelegator{
	name:    "Plaintext",
	parse:   ParsePlaintextValue,
	compare: ConstantValueCompare,
}

type PlaintextValue string

func ParsePlaintextValue(s string) (Value, error) {
	return PlaintextValue(s), nil
}

// IsPlaintextValue checks if the provided value
// matches the criteria for a plaintext password
// (i.e. starts with an alphanumeric character).
//
// this makes plaintext passwords safe to use as
// fallback codec, as other characters might be
// used as indicator for other codecs.
func IsPlaintextValue(s string) bool {
	initial, _ := utf8.DecodeRuneInString(s)
	if initial == utf8.RuneError {
		return false
	}

	return (unicode.IsNumber(initial) || unicode.IsLetter(initial))
}

func (p PlaintextValue) Generate(s string) (string, error) {
	return s, nil
}

func (p PlaintextValue) String() string {
	return string(p)
}
