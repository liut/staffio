package admin

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	staffio "github.com/liut/staffio/client"
)

const (
	sKeyUser = "user"
)

var (
	UserLifetime int64 = 3600
	Guest              = &User{}

	AdminPath = "/admin/"
	LoginPath = "/auth/login"

	LoginHandler = gin.WrapF(staffio.LoginHandler)
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	LastHit int64  `json:"-"`
}

func (u *User) IsExpired() bool {
	if UserLifetime == 0 {
		return false
	}
	return u.LastHit+UserLifetime < time.Now().Unix()
}

func (u *User) needRefresh() bool {
	return time.Now().Unix()-u.LastHit < UserLifetime/2
}

func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}

func AuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := ginSession(c)
		if user, ok := sess.Get(sKeyUser).(*User); ok {
			if !user.IsExpired() {
				if user.needRefresh() {
					user.Refresh()
					sess.Set(sKeyUser, user)
					smgr.Save(sess, c.Writer)
				}
				c.Set(sKeyUser, user)
				c.Next()
				return
			}
		}

		if redirect {
			c.Redirect(http.StatusFound, LoginPath)
			c.Abort()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func UserWithContext(c *gin.Context) (user *User, ok bool) {
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
	hf := func(w http.ResponseWriter, r *http.Request) {
		it, err := staffio.AuthRequestWithRole(r, roleName...)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("auth with role %v ERR %s", roleName, err)
			return
		}

		sess := loadSession(r)
		user := &User{
			Uid:  it.Me.Uid,
			Name: it.Me.Nickname,
		}
		user.Refresh()
		sess.Set(sKeyUser, user)
		smgr.Save(sess, w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Refresh", fmt.Sprintf("0; %s", AdminPath))
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Login OK. Please waiting, ok click <a href=" + AdminPath + ">here</a> to go back"))
		return
	}
	return gin.WrapH(staffio.AuthCodeCallbackWrap(http.HandlerFunc(hf)))
}

func HandlerShowMe(c *gin.Context) {
	user, ok := UserWithContext(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"me": user,
	})
}
