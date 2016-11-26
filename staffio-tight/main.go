package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rakyll/statik/fs"
	"github.com/wealthworks/go-utils/reaper"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
	_ "lcgc/platform/staffio/staffio-tight/statik"
	"lcgc/platform/staffio/web"
)

func main() {
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))
	router := web.MainRouter()

	statikFS, se := fs.New()
	if se != nil {
		log.Fatalf(se.Error())
	}

	ss := http.FileServer(statikFS)
	router.Path("/favicon.ico").Handler(ss).Methods("GET", "HEAD")
	router.Path("/robots.txt").Handler(ss).Methods("GET", "HEAD")

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)
	err := http.ListenAndServe(Settings.HttpListen, router) // Start the server!
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
