package web

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	zlog "github.com/liut/staffio/pkg/log"
)

func TestMain(m *testing.M) {

	lgr, _ := zap.NewDevelopment()
	defer func() {
		_ = lgr.Sync() // flushes buffer, if any
	}()
	sugar := lgr.Sugar()
	zlog.SetLogger(sugar)

	os.Exit(m.Run())
}

func TestUrlFor(t *testing.T) {

	loginUrl := UrlFor("login")
	assert.Equal(t, "/login", loginUrl)
}

func TestSchemaClient(t *testing.T) {
	s := `{"id":1,"name":"test2","redirect_uri":"http://localhost:3001"}`
	var c clientParam
	err := json.Unmarshal([]byte(s), &c)
	assert.NoError(t, err)
	assert.Equal(t, 1, c.ID)
}
