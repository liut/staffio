package client

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type ctxKey int

const (
	TokenKey ctxKey = 0
)

func setCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// AuthCodeCallback is a middleware that injects a InfoToken into the context of each request
func AuthCodeCallback(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// verify state value.
		stateCookie, err := r.Cookie(cKeyState)
		if err != nil || stateCookie.Value != r.FormValue("state") {
			log.Printf("Invalid state:\n%s\n%s", stateCookie.Value, r.FormValue("state"))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid state: " + stateCookie.Value))
			return
		}

		tok, err := conf.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil {
			log.Printf("oauth2 exchange ERR %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		deleteCookie(w, cKeyState)
		if RememberToken {
			setCookie(w, cKeyToken, tok.AccessToken)
		}
		if RememberUser {
			if uid, ok := tok.Extra("uid").(string); ok {
				setCookie(w, cKeyUser, uid)
			}
		}

		log.Printf("exchanged token: %s", tok)

		token, err := requestInfoToken(tok)
		if err != nil {
			log.Printf("requestInfoToken err %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, TokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// TokenFromContext returns a InfoToken from the given context if one is present.
// Returns nil if a InfoToken cannot be found.
func TokenFromContext(ctx context.Context) *InfoToken {
	if ctx == nil {
		return nil
	}
	if tok, ok := ctx.Value(TokenKey).(*InfoToken); ok {
		return tok
	}
	return nil
}
