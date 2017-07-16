package web

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	. "lcgc/platform/staffio/pkg/settings"
)

var (
	CookieName = "_user"
)

const (
	kAuthUser = "user"
	LoginPath = "/login"
)

func AuthUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := UserFromRequest(c.Request)
		if err != nil {
			markReferer(c)
			c.Redirect(302, LoginPath)
			c.Abort()
			return
		}
		// log.Printf("got user %q", user.Uid)
		c.Set(kAuthUser, user)
	}
}

func (s *server) AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exist := c.Get(kAuthUser)
		if !exist {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Abort()
			return
		}
		user := v.(*User)
		if !s.IsKeeper(user.Uid) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Abort()
		}
	}
}

func UserWithContext(c *gin.Context) (user *User) {
	if v, ok := c.Get(kAuthUser); ok {
		user = v.(*User)
	}
	if user == nil {
		panic("user not found in request")
	}

	return
}

func UserFromRequest(r *http.Request) (user *User, err error) {
	var cookie *http.Cookie
	cookie, err = r.Cookie(CookieName)
	if err != nil {
		log.Printf("cookie %q ERR %s", CookieName, err)
		return
	}
	var b []byte
	b, err = base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Printf("base64decode %q ERR %s", cookie.Value, err)
		return
	}
	// log.Printf("got encrypted %s", b)
	user = new(User)
	err = user.Decode(b)
	if err != nil {
		log.Printf("decode msgpack ERR %s", err)
	}
	// log.Printf("got user %v", user)
	return
}

func (user *User) toResponse(w http.ResponseWriter) error {
	b, err := user.Encode()
	if err != nil {
		log.Printf("marshal msgpack user ERR: %s", err)
		return err
	}
	value := base64.URLEncoding.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  value,
		MaxAge: Settings.UserLifetime,
		Path:   "/",
		// Domain:   "",
		// Secure:   false,
		HttpOnly: true,
	})
	return nil
}
