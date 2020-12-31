package password

import (
	"fmt"
)

type Generator interface {
	Generate(string) (string, error)
}

type Value interface {
	Generator
	fmt.Stringer
}

type Password struct {
	Value
	Codec
}

func ParsePassword(s string) (pass *Password, err error) {
	var hash string
	pass = &Password{}

	pass.Codec, hash, err = FindCodec(s)
	if err != nil {
		return
	}

	pass.Value, err = pass.Codec.ParseValue(hash)
	if err != nil {
		return
	}

	return
}

func NewPlaintextPassword(s string) *Password {
	password := &Password{
		Codec: plaintextCodec,
		Value: PlaintextValue(s),
	}

	return password
}

func (p *Password) Compare(plain string) bool {
	return p.Codec.Compare(plain, p.Value)
}
