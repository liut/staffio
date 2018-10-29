package admin

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	staffio "github.com/liut/staffio/client"
)

const (
	sKeyUser = "user"
	KeyOper  = "oper"
)

type User = staffio.User

var (
	LoginHandler = gin.WrapF(staffio.LoginHandler)
	SetLoginPath = staffio.SetLoginPath
	SetAdminPath = staffio.SetAdminPath
)

func AuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := staffio.UserFromRequest(c.Request)
		if err != nil {
			if redirect {
				c.Redirect(http.StatusFound, staffio.LoginPath)
				c.Abort()
				return
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(c.Request.Context(), KeyOper, user)
		c.Request = c.Request.WithContext(ctx)
		c.Set(sKeyUser, user)
		c.Next()
	}
}

func UserFromContext(c *gin.Context) (user *User, ok bool) {
	v, ok := c.Get(sKeyUser)
	if ok {
		user = v.(*User)
	}
	if user == nil {
		log.Print("user not found in request")
	}

	return
}

// AuthCodeCallback Handler for Check auth with role[s] when auth-code callback
func AuthCodeCallback(roleName ...string) gin.HandlerFunc {
	return gin.WrapH(staffio.AuthCodeCallback(roleName...))
}

func HandlerShowMe(c *gin.Context) {
	user, ok := UserFromContext(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, _ := user.Encode()
	c.JSON(http.StatusOK, gin.H{
		"me":    user,
		"token": token,
	})
}
