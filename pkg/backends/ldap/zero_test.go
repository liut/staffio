package ldap

import (
	"testing"
)

func TestSplitDC(t *testing.T) {
	base := "dc=example,dc=org"
	dc1 := splitDC(base)

	if dc1 != "example" {
		t.Errorf("mismatch %q and %q", dc1, "example")
	}
}
