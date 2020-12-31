package password

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptParamDelimiter = "$"
	BcryptIndicator      = "$2y$"
)

var ErrMalformedBcrypt = errors.New("malformed Bcrypt value")
var bcryptCodec = &CodecDelegator{
	name:      "Bcrypt",
	indicator: BcryptIndicator,
	parser:    ParseBcryptValue,
}

type BcryptValue struct {
	Cost int
	Hash string
}

func ParseBcryptValue(s string) (Value, error) {
	params := strings.Split(s, BcryptParamDelimiter)

	if len(params) != 2 {
		return nil, ErrMalformedBcrypt
	}

	cost, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	value := &BcryptValue{
		Cost: cost,
		Hash: params[1],
	}
	return value, nil
}

func (p *BcryptValue) Generate(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), p.Cost)
	if err != nil {
		return "", err
	}

	return p.assemble(string(hash)), nil
}

func (p *BcryptValue) String() string {
	return p.assemble(p.Hash)
}

func (p *BcryptValue) assemble(hash string) string {
	cost := strconv.Itoa(p.Cost)

	return fmt.Sprint(BcryptIndicator, cost, BcryptParamDelimiter, hash)
}
