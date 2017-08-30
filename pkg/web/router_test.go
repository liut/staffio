package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlFor(t *testing.T) {

	loginUrl := UrlFor("login")
	assert.Equal(t, "/login", loginUrl)
}
