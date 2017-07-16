package web

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {

	stamp := time.Now().Unix()
	user := &User{"test", "name", stamp, 1}
	b, err := user.Encode()
	assert.Nil(t, err)
	assert.True(t, len(b) > 0)
	var out = new(User)
	err = out.Decode(b)
	assert.Nil(t, err)
	assert.NotEmpty(t, out.Uid, out.Name)
	assert.NotZero(t, out.LastHit)
	assert.Equal(t, user.Uid, out.Uid)
}
