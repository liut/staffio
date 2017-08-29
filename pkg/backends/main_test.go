package backends

import (
	"log"
	"testing"

	"github.com/liut/staffio/pkg/settings"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	settings.Parse()
	m.Run()
}
