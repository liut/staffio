package types

import (
	"testing"
)

func TestIID(t *testing.T) {
	gums := []struct {
		v uint64
		s string
	}{
		{149495437762496513, "14vzpk09yxoh"},
		{149497847983638530, "14w0kb8xep6q"},
	}
	for _, i := range gums {
		checkEqual(t, i.s, IID(i.v).String())
		var id = new(IID)
		err := id.UnmarshalText([]byte(i.s))
		if err != nil {
			t.Fail()
		}
		checkEqual(t, i.v, uint64(*id))
	}
}

func checkEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("Not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", expected, actual)
	}

}
