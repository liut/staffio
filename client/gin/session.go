package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-osin/session"
)

var (
	smgr       session.Manager
	sessionKey = "gin-session"

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

func loadSession(r *http.Request) session.Session {
	sess := smgr.Load(r)
	if sess == nil {
		sess = session.NewSession()
	}
	return sess
}

func ginSession(c *gin.Context) session.Session {
	if sess, ok := c.Get(sessionKey); ok {
		return sess.(session.Session)
	}
	sess := loadSession(c.Request)
	c.Set(sessionKey, sess)
	return sess
}
