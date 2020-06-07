package backends

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenClient(t *testing.T) {
	names := []string{"test 1", "test 2"}
	for _, name := range names {
		c := GenNewClient(name, "")
		assert.NotNil(t, c)
		t.Logf("new client %s", c)
	}
}
