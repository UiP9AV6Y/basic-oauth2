package password

import (
	"crypto/subtle"
)

type ValueParser func(string) (Value, error)

type CodecDelegator struct {
	name      string
	indicator string
	parser    ValueParser
}

func (d *CodecDelegator) Name() string {
	return d.name
}

func (d *CodecDelegator) Indicator() string {
	return d.indicator
}

func (d *CodecDelegator) ParseValue(s string) (Value, error) {
	return d.parser(s)
}

func (d *CodecDelegator) Compare(s string, v Value) bool {
	input, err := v.Generate(s)
	if err != nil {
		return false
	}

	x := []byte(v.String())
	y := []byte(input)

	return (subtle.ConstantTimeCompare(x, y) == 1)
}
