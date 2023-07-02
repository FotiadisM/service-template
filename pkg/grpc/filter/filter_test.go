package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFullMethodName(t *testing.T) {
	s, m := ParseFullMethodName("/auth.v1.AuthService/Register")
	assert.Equal(t, "auth.v1.AuthService", s)
	assert.Equal(t, "Register", m)
}
