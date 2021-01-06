package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstantValueCompare(t *testing.T) {
	value := PlaintextValue("password")
	gotMatch := ConstantValueCompare("password", value)
	gotMismatch := ConstantValueCompare("username", value)

	assert.True(t, gotMatch, "password/password")
	assert.False(t, gotMismatch, "username/password")
}
