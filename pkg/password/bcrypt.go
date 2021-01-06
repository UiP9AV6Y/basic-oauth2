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
	parse:     ParseBcryptValue,
	compare:   BcryptValueCompare,
	indicator: BcryptIndicator,
}

// Password encoded using the Bcrypt hashing algorithm
type BcryptValue struct {
	Cost int
	Hash string
}

// ParseBcryptValue parses the input value for
// Bcrypt hashing parameters. the value must not
// be prefixed with the indicator but merely consists
// of the encoded parameters.
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

// BcryptValueCompare is a wrapper around
// bcrypt.CompareHashAndPassword()
func BcryptValueCompare(s string, v Value) bool {
	x := []byte(v.String())
	y := []byte(s)

	return bcrypt.CompareHashAndPassword(x, y) == nil
}

// Generate uses the input value to generate a password using
// the same parameters as this instance. The output is compatible
// with BcryptValueCompare() but might not be accepted by
// ParsePassword()/FindCodec() due to different versions
// of the bcrypt algorithm being used.
func (p *BcryptValue) Generate(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), p.Cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// String encodes the parameters of this value.
func (p *BcryptValue) String() string {
	cost := strconv.Itoa(p.Cost)

	return fmt.Sprint(BcryptIndicator, cost, BcryptParamDelimiter, p.Hash)
}
