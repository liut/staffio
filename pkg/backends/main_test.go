package backends

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/liut/staffio/pkg/log"
)

var (
	svc Servicer
)

func TestMain(m *testing.M) {
	_logger, _ := zap.NewDevelopment()
	defer _logger.Sync() // flushes buffer, if any
	sugar := _logger.Sugar()
	log.SetLogger(sugar)
	SetDSN(os.Getenv("STAFFIO_BACKEND_DSN"))
	svc = NewService()
	svc.Ready()
	m.Run()
}

func TestWatching(t *testing.T) {
	uid := "eagle"
	data := svc.Watch().Gets(uid)
	assert.Empty(t, data)
	target := "john"

	err := svc.Watch().Watch(uid, target)
	assert.NoError(t, err)
	data = svc.Watch().Gets(uid)
	assert.NotEmpty(t, data)
	t.Logf("watching of %s: %v", uid, data)

	err = svc.Watch().Unwatch(uid, target)
	assert.NoError(t, err)
}
