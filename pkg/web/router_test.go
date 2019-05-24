package web

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlFor(t *testing.T) {

	loginUrl := UrlFor("login")
	assert.Equal(t, "/login", loginUrl)
}

func TestSchemaClient(t *testing.T) {
	s := `{"id":1,"name":"test2","redirect_uri":"http://localhost:3001"}`
	var c clientParam
	err := json.Unmarshal([]byte(s), &c)
	assert.NoError(t, err)
	assert.Equal(t, 1, c.ID)
}
