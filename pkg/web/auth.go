package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	auth "github.com/liut/simpauth"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

var (
	authzr auth.Authorizer
)

func init() {
	authzr = auth.New(auth.WithURI("/login"), auth.WithCookie(
		settings.Current.CookieName,
		settings.Current.CookiePath,
		settings.Current.CookieDomain,
	), auth.WithMaxAge(settings.Current.CookieMaxAge))
}

// consts
const (
	kAuthUser = "user"
)

type User = auth.User

func UserFromStaff(staff *models.Staff) *auth.User {
	return &auth.User{
		Avatar: staff.AvatarURI(),
		UID:    staff.UID,
		Name:   staff.GetName(),
	}
}

func AuthUserMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authzr.UserFromRequest(c.Request)
		if err != nil {
			logger().Infow("user from request", "err", err)
			if redirect {
				markReferer(c)
				c.Redirect(302, UrlFor("login"))
				c.Abort()
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				// apiError(c, ERROR_INTERNAL, err)
			}
			return
		}
		// log.Printf("got user %q", user.UID)
		c.Set(kAuthUser, user)
		c.Next()
		user.Refresh()
		_ = user.Signin(c.Writer)
	}
}

func (s *server) authGroup(name ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exist := c.Get(kAuthUser)
		if !exist {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Abort()
			return
		}
		user := v.(*User)
		if !s.InGroupAny(user.UID, name...) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Abort()
			return
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

func signinStaffGin(c *gin.Context, staff *models.Staff) {
	user := UserFromStaff(staff)
	user.Refresh()
	logger().Debugw("login ok", "user", user)
	sess := ginSession(c)
	sess.Set(kAuthUser, user)
	_ = authzr.Signin(user, c.Writer)
	SessionSave(sess, c.Writer)
}
