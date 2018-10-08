package client

import (
	"log"
	"net/http"

	"github.com/go-osin/session"
)

// TODO: deprecated with cookie

const (
	SessKeyUser  = "user"
	SessKeyToken = "token"
)

var (
	SessionIDCookieName = "_sess"
)

func init() {
	SetupSessionStore(session.NewInMemStore())
}

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

func UserFromSession(sess session.Session) (u *User, ok bool) {
	u, ok = sess.Get(SessKeyUser).(*User)
	return
}

func (user *User) SaveToSession(sess session.Session, w http.ResponseWriter, force bool) {
	if force || user.NeedRefresh() {
		user.Refresh()
		sess.Set(SessKeyUser, user)
		SessionSave(sess, w)
		log.Printf("saved user %v to session", user)
	}
}
