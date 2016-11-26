package web

import (
	"log"

	"github.com/liut/pgstore"

	. "lcgc/platform/staffio/settings"
)

var (
	store *pgstore.PGStore
)

func sessionInit() {
	var err error
	// store = sessions.NewCookieStore([]byte(Settings.Session.Name))
	store, err = pgstore.NewPGStore(Settings.Backend.DSN, []byte(Settings.Session.Secret))
	if err != nil {
		log.Fatal(err)
	}
	store.MaxAge(Settings.Session.MaxAge)
	// store.Options.MaxAge = Settings.Session.MaxAge
	store.Options.Domain = Settings.Session.Domain
	store.Options.HttpOnly = true
}
