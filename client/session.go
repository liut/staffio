package client

import (
	"net/http"

	"github.com/go-osin/session"
)

const (
	SessKeyUser = "user"
)

var (
	SessionIDCookieName = "_sess"
)

func SetupSessionStore(store session.Store) {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(store, &session.CookieMngrOptions{
		SessIDCookieName: SessionIDCookieName,
		AllowHTTP:        true,
	})
}

func SessionLoad(r *http.Request) session.Session {
	sess := session.Global.Load(r)
	if sess == nil {
		sess = session.NewSession()
	}
	return sess
}

func SessionSave(sess session.Session, w http.ResponseWriter) {
	session.Global.Save(sess, w)
}
