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

// AuthCodeCallbackWrap is a middleware that injects a InfoToken with roles into the context of each request
func AuthCodeCallbackWrap(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// verify state value.
		stateCookie, err := r.Cookie(cKeyState)
		if err != nil {
			log.Printf("cookie not found: %s", cKeyState)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if stateCookie.Value != r.FormValue("state") {
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

		log.Printf("exchanged token: %s", tok)

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
