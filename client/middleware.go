package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	ErrNoToken = errors.New("oauth2 token not found")
	ErrNoRole  = errors.New("the user not in special roles")

	AdminPath = "/admin/"
	LoginPath = "/auth/login"
)

type ctxKey int

const (
	TokenKey ctxKey = iota
	UserKey
)

func SetLoginPath(path string) {
	LoginPath = path
}

func SetAdminPath(path string) {
	AdminPath = path
}

// AuthMiddleware
func AuthMiddleware(redirect bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			user, err := UserFromRequest(r)
			if err != nil {
				if redirect {
					http.Redirect(w, r, LoginPath, http.StatusFound)
					return
				}
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hf)
	}
}

// UserFromContext
func UserFromContext(ctx context.Context) *User {
	if ctx == nil {
		return nil
	}
	if obj, ok := ctx.Value(UserKey).(*User); ok {
		return obj
	}
	return nil
}

// AuthCodeCallback Handler for Check auth with role[s] when auth-code callback
func AuthCodeCallback(roleName ...string) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		it, err := AuthRequestWithRole(r, roleName...)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("auth with role %v ERR %s", roleName, err)
			return
		}

		user := &User{
			Uid:  it.Me.Uid,
			Name: it.Me.Nickname,
		}
		user.Refresh()
		Signin(user, w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Refresh", fmt.Sprintf("0; %s", AdminPath))
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Login OK. Please waiting, ok click <a href=" + AdminPath + ">here</a> to go back"))
		return
	}
	return AuthCodeCallbackWrap(http.HandlerFunc(hf))
}

// AuthCodeCallbackWrap is a middleware that injects a InfoToken with roles into the context of callback request
func AuthCodeCallbackWrap(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// verify state value.
		sess := SessionLoad(r)
		var state string
		v := sess.Get(cKeyState)
		if v == nil {
			log.Printf("last state not found: %s", cKeyState)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		state = v.(string)
		if state != r.FormValue("state") {
			log.Printf("Invalid state at %s:\n%s\n%s", r.RequestURI, state, r.FormValue("state"))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid state: " + state))
			return
		}

		tok, err := conf.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil {
			log.Printf("oauth2 exchange ERR %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// log.Printf("exchanged token: %s", tok)
		sess.Set(cKeyToken, tok.AccessToken)

		ctx := r.Context()
		ctx = context.WithValue(ctx, TokenKey, tok)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// UidFromToken extract uid from oauth2.Token
func UidFromToken(tok *oauth2.Token) string {
	if uid, ok := tok.Extra("uid").(string); ok {
		return uid
	}
	return ""
}

// TokenFromContext returns a oauth2.Token from the given context if one is present.
// Returns nil if a oauth2.Token cannot be found.
func TokenFromContext(ctx context.Context) *oauth2.Token {
	if ctx == nil {
		return nil
	}
	if tok, ok := ctx.Value(TokenKey).(*oauth2.Token); ok {
		return tok
	}
	return nil
}

// AuthRequestWithRole called in AuthCallback
func AuthRequestWithRole(r *http.Request, role ...string) (it *InfoToken, err error) {
	tok := TokenFromContext(r.Context())
	if tok == nil {
		err = ErrNoToken
		return
	}
	it, err = RequestInfoToken(tok, role...)
	if err != nil {
		return
	}
	for _, rn := range role {
		if !it.Roles.Has(rn) {
			err = ErrNoRole
			break
		}
	}

	return
}
