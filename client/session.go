package client

import (
	"net/http"

	"github.com/go-osin/session"
)

const (
	SessKeyUser = "user"
)

var (
	smgr session.Manager

	SessionIDCookieName = "_sess"
)

func init() {
	SetupSessionStore(session.NewInMemStore())
}

func SetupSessionStore(store session.Store) {
	smgr = session.NewCookieManagerOptions(store, &session.CookieMngrOptions{
		SessIDCookieName: SessionIDCookieName,
		AllowHTTP:        true,
	})
}

func SessionFromRequest(r *http.Request) session.Session {
	sess := smgr.Load(r)
	if sess == nil {
		sess = session.NewSession()
	}
	return sess
}

func SessionSave(sess session.Session, w http.ResponseWriter) {
	smgr.Save(sess, w)
}
