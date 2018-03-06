package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-osin/session"
)

var (
	smgr       session.Manager
	sessionKey = "gin-session"
)

func init() {
	smgr = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{
		SessIDCookieName: "st_sess",
		AllowHTTP:        true,
	})
}

func (s *server) loadSession(r *http.Request) session.Session {
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
	sess := svr.loadSession(c.Request)
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
