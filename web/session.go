package web

import (
	"log"
	"time"

	"github.com/liut/pgstore"

	. "lcgc/platform/staffio/settings"
)

var (
	store *pgstore.PGStore
	quit  chan<- struct{}
	done  <-chan struct{}
)

func sessionStart() {
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

	quit, done = store.Cleanup(time.Minute * 5)
}

func sessionStop() {
	store.StopCleanup(quit, done)
}
