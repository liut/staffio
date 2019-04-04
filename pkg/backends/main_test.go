package backends

import (
	"log"
	"testing"

	"github.com/liut/staffio/pkg/settings"
)

var (
	svc Servicer
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	settings.Parse()
	svc = NewService()
	svc.Ready()
	m.Run()
}
