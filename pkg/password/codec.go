package password

import (
	"errors"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
)

var ErrUnsupportedCodec = errors.New("unsupported password codec")

var passwordCodecs = []Codec{
	bcryptCodec,
	shaCodec,
}

type Codec interface {
	Name() string
	Indicator() string
	Compare(string, Value) bool
	ParseValue(string) (Value, error)
}

func FindCodec(s string) (codec Codec, params string, err error) {
	for _, c := range passwordCodecs {
		if p, ok := utils.SplitPrefix(s, c.Indicator()); ok {
			codec = c
			params = p
			return
		}
	}

	if IsPlaintextValue(s) {
		codec = plaintextCodec
		params = s
		return
	}

	err = ErrUnsupportedCodec
	return
}
