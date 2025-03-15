package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var (
		code        = "a01"
		secret      = "secret"
		redirectURI = "http://localhost"
		meta        = defaultClientMeta
	)
	meta.Name = "test"
	c := NewClient(code, secret, redirectURI)
	c.Meta = meta

	assert.Equal(t, code, c.GetId())
	assert.Equal(t, secret, c.GetSecret())
	assert.Equal(t, redirectURI, c.GetRedirectUri())
	assert.Equal(t, meta, c.Meta)
	assert.Equal(t, "test", c.GetName())

}

func TestJSONKV(t *testing.T) {
	m := JSONKV{"name": "eagle"}
	v, ok := m.Get("name")
	assert.True(t, ok)
	assert.Equal(t, v, "eagle")
}
