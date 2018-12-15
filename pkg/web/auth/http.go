package auth

import (
	"fmt"
	"log"
	"net/http"
)

var (
	CookieName   = "_user"
	CookiePath   = "/"
	CookieMaxAge = 3600
)

// UserFromRequest get user from cookie
func UserFromRequest(r *http.Request) (user *User, err error) {
	var cookie *http.Cookie
	cookie, err = r.Cookie(CookieName)
	if err != nil {
		log.Printf("cookie %q ERR %s", CookieName, err)
		return
	}
	user = new(User)
	err = user.Decode(cookie.Value)
	if err != nil {
		log.Printf("decode user ERR %s", err)
		return
	}
	if user.IsExpired() {
		err = fmt.Errorf("user %s is expired", user.UID)
	}
	// log.Printf("got user %v", user)
	return
}

// Signin call Signin for login
func (user *User) Signin(w http.ResponseWriter) error {
	return Signin(user, w)
}

type encoder interface {
	Encode() (string, error)
}

// Signin write user encoded string into cookie
func Signin(user encoder, w http.ResponseWriter) error {
	value, err := user.Encode()
	if err != nil {
		log.Printf("encode user ERR: %s", err)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    value,
		MaxAge:   CookieMaxAge,
		Path:     CookiePath,
		HttpOnly: true,
	})
	return nil
}

// Signout setcookie with empty
func Signout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     CookiePath,
		HttpOnly: true,
	})
}
