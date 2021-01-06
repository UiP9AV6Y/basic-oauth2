package password

import (
	"crypto/sha1"
	"encoding/base64"
)

const (
	SHAIndicator = "{SHA}"
)

var shaCodec = &CodecDelegator{
	name:      "SHA",
	parse:     ParseSHAValue,
	compare:   ConstantValueCompare,
	indicator: SHAIndicator,
}

type SHAValue string

func ParseSHAValue(s string) (Value, error) {
	return SHAValue(s), nil
}

func (p SHAValue) Generate(s string) (string, error) {
	raw := sha1.Sum([]byte(s))
	hash := base64.StdEncoding.EncodeToString(raw[:])

	return (SHAIndicator + hash), nil
}

func (p SHAValue) String() string {
	return SHAIndicator + string(p)
}
