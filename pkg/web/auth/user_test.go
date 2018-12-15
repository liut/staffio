package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {

	stamp := time.Now().Unix()
	user := &User{UID: "test", Name: "name", LastHit: stamp}
	b, err := user.Encode()
	assert.Nil(t, err)
	assert.True(t, len(b) > 0)
	var out = new(User)
	err = out.Decode(b)
	assert.Nil(t, err)
	assert.NotEmpty(t, out.UID, out.Name)
	assert.NotZero(t, out.LastHit)
	assert.Equal(t, user.UID, out.UID)
}

func TestUserEncode(t *testing.T) {
	var user = &User{
		UID:  "test",
		Name: "testname",
	}
	s, err := user.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encoded %q", s)
	var other = new(User)
	err = other.Decode(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user %q", other.UID)
}

func TestUserDecode(t *testing.T) {
	// s := "g6F1pmxpdXRhb6FupWxpw7p0oWjSW7r67w"
	s := "g6F1pHRlc3Shbqh0ZXN0bmFtZaFoAA"
	var user = new(User)
	err := user.Decode(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user %q", user.UID)
}
