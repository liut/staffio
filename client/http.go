package client

import (
	"fmt"
	"log"
	"net/http"
)

var (
	CookieName   = "_user"
	CookiePath   = "/"
	CookieMaxAge = 3600

	// ErrInvalidToken = errors.New("invalid token or expired")
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
		log.Printf("decode cookie ERR %s", err)
		return
	}
	if user.IsExpired() {
		err = fmt.Errorf("user %q from %q is expired", user.Uid, cookie.Value)
		log.Print(err)
	}

	return
}

type encoder interface {
	Encode() (string, error)
}

// Signin write user encoded string into cookie
func Signin(user encoder, w http.ResponseWriter) error {
	value, err := user.Encode()
	if err != nil {
		log.Printf("call user.Encode() ERR: %s", err)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    value,
		MaxAge:   CookieMaxAge,
		Path:     CookiePath,
		HttpOnly: true,
	})
	log.Printf("signin user %v, %q", user, value)
	return nil
}

// Signout clear cookie for user
func Signout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     CookiePath,
		HttpOnly: true,
	})
}
