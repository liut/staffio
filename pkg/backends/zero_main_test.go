package backends

import (
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/liut/staffio/pkg/log"
)

var (
	svc   Servicer
	store OSINStore
)

func TestMain(m *testing.M) {
	_logger, _ := zap.NewDevelopment()
	defer func() {
		_ = _logger.Sync() // flushes buffer, if any
	}()
	sugar := _logger.Sugar()
	log.SetLogger(sugar)

	SetDSN(envOr("STAFFIO_BACKEND_TEST_DSN", "postgres://staffio@localhost/staffiotest?sslmode=disable"))

	db := getDb()
	logger().Infow("cleaning schemas")
	db.Exec("DROP SCHEMA IF EXISTS staffio CASCADE;") //nolint

	schemas := []string{
		"staffio_0_schema.sql",
		"staffio_1_schema.sql",
		"staffio_2_schema_cas.sql",
		"staffio_2_schema_session.sql",
		"staffio_3_schema_team.sql",
		"staffio_3_schema_weekly.sql",
		"staffio_4_schema_content.sql",
		"staffio_5_init.sql",
	}
	for _, fn := range schemas {
		if _e := execSQLfile(db, "../../database/"+fn); _e != nil {
			panic(_e)
		}
	}

	defer func() {
		logger().Infow("test done")
	}()

	svc = NewService()
	store = svc.OSIN()
	svc.Ready() //nolint
	code := m.Run()
	os.Exit(code)
}

func loadSQLs(name string) string {
	content, err := os.ReadFile(name)
	if err != nil {
		logger().Fatalw("loadSQL fail", "err", err)
	}
	return string(content)
}

func execSQLfile(db dber, name string) error {
	query := strings.TrimSpace(loadSQLs(name))
	if query == "" {
		return nil
	}
	_, err := db.Exec(query)
	if err != nil {
		logger().Warnw("exec sql fail", "name", name, "query", query[:32], "err", err)
		return err
	}
	return nil
}

func TestCleanup(t *testing.T) {
	err := Cleanup()
	assert.NoError(t, err)
}

func TestWatching(t *testing.T) {
	uid := "eagle"
	assert.NotNil(t, svc.Watch().Gets(uid).UIDs())
	data := svc.Watch().Gets(uid)
	assert.Empty(t, data)
	target := "john"

	err := svc.Watch().Watch(uid, target)
	assert.NoError(t, err)
	data = svc.Watch().Gets(uid)
	if assert.NotEmpty(t, data) {
		assert.Equal(t, target, data[0].UID)
	}

	t.Logf("watching of %s: %v", uid, data)

	err = svc.Watch().Unwatch(uid, target)
	assert.NoError(t, err)
}
