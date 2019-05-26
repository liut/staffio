package backends

import (
	"os"
	"testing"
)

var (
	svc Servicer
)

func TestMain(m *testing.M) {
	SetDSN(os.Getenv("STAFFIO_BACKEND_DSN"))
	svc = NewService()
	svc.Ready()
	m.Run()
}
