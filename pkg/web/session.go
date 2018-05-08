package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-osin/session"
)

var (
	once       sync.Once
	smgr       session.Manager
	sessionKey = "gin-session"
)

func SetupSessionStore(store session.Store) {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(store, &session.CookieMngrOptions{
		SessIDCookieName: "st_sess",
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

func ginSession(c *gin.Context) session.Session {
	if sess, ok := c.Get(sessionKey); ok {
		return sess.(session.Session)
	}
	sess := SessionLoad(c.Request)
	c.Set(sessionKey, sess)
	return sess
}

func (s *server) sessionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := ginSession(c)
		defer func() {
			if sess.Changed() {
				smgr.Save(sess, c.Writer)
			}
		}()
		c.Next()
	}
}
