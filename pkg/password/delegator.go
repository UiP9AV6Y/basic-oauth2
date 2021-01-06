package password

import (
	"crypto/subtle"
)

type ValueParser func(string) (Value, error)

type ValueComparator func(string, Value) bool

func ConstantValueCompare(s string, v Value) bool {
	input, err := v.Generate(s)
	if err != nil {
		return false
	}

	x := []byte(v.String())
	y := []byte(input)

	return (subtle.ConstantTimeCompare(x, y) == 1)
}

type CodecDelegator struct {
	name      string
	indicator string
	compare   ValueComparator
	parse     ValueParser
}

func (d *CodecDelegator) Name() string {
	return d.name
}

func (d *CodecDelegator) Indicator() string {
	return d.indicator
}

func (d *CodecDelegator) ParseValue(s string) (Value, error) {
	return d.parse(s)
}

func (d *CodecDelegator) Compare(s string, v Value) bool {
	return d.compare(s, v)
}
