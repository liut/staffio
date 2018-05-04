package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-osin/session"

	staffio "github.com/liut/staffio/client"
)

var (
	sessionKey = "gin-session"
)

func ginSession(c *gin.Context) session.Session {
	if sess, ok := c.Get(sessionKey); ok {
		return sess.(session.Session)
	}
	sess := staffio.SessionFromRequest(c.Request)
	c.Set(sessionKey, sess)
	return sess
}
