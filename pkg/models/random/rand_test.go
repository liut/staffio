package random

import (
	"testing"
)

func TestGenString(t *testing.T) {
	s := GenString(32)
	t.Logf("generated randam string %q", s)
}
