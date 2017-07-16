package web

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	gsess "github.com/gorilla/sessions"
	"github.com/liut/pgstore"

	. "lcgc/platform/staffio/pkg/settings"
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
	// store = sessions.NewCookieStore([]byte(Settings.Session.Name))
	store, err := pgstore.NewPGStore(Settings.Backend.DSN, []byte(Settings.Session.Secret))
	if err != nil {
		log.Fatal(err)
	}
	store.MaxAge(Settings.Session.MaxAge)
	// store.Options.MaxAge = Settings.Session.MaxAge
	store.Options.Domain = Settings.Session.Domain
	store.Options.HttpOnly = true

	return &pgStore{store}
}

func ginSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}
