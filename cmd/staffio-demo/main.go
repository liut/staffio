package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/liut/keeper/utils/reaper"

	"github.com/liut/staffio/pkg/backends"
	config "github.com/liut/staffio/pkg/settings"
	"github.com/liut/staffio/pkg/web"
)

const (
	readTimeout  time.Duration = 10 * time.Second
	writeTimeout               = 15 * time.Second
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	settings := config.Current

	cfg := web.Config{
		Root:    ".",
		FS:      "local",
		BaseURI: settings.BaseURL,
	}
	ws := web.New(cfg)
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	srv := &http.Server{
		Addr:         settings.HTTPListen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      ws,
	}

	d := &demo{
		base: settings.BaseURL,
	}
	d.strap(ws)

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", config.Version(), settings.HTTPListen, settings.Root)
	err := srv.ListenAndServe() // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}
}
