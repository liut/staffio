package client

import (
	"testing"
)

func TestUserMarshal(t *testing.T) {
	var user = &User{
		Uid:  "test",
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
	t.Logf("user %q", other.Uid)
}

func TestUserUnmarshal(t *testing.T) {
	// s := "g6F1pmxpdXRhb6FupWxpw7p0oWjSW7r67w"
	s := "g6F1pHRlc3Shbqh0ZXN0bmFtZaFoAA"
	var user = new(User)
	err := user.Decode(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user %q", user.Uid)
}
