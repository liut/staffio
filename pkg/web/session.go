package web

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	gsess "github.com/gorilla/sessions"
	"github.com/liut/pgstore"

	"github.com/liut/staffio/pkg/settings"
)

type pgStore struct {
	*pgstore.PGStore
}

func (c *pgStore) Options(options sessions.Options) {
	c.PGStore.Options = &gsess.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func sessionStore() sessions.Store {
	var err error
	// store = sessions.NewCookieStore([]byte(settings.Session.Name))
	store, err := pgstore.NewPGStore(settings.Backend.DSN, []byte(settings.Session.Secret))
	if err != nil {
		log.Fatal(err)
	}
	store.MaxAge(settings.Session.MaxAge)
	// store.Options.MaxAge = settings.Session.MaxAge
	store.Options.Domain = settings.Session.Domain
	store.Options.HttpOnly = true

	return &pgStore{store}
}

func ginSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}
