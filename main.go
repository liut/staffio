package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wealthworks/go-utils/reaper"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
	"lcgc/platform/staffio/web"
	"lcgc/platform/staffio/webfatso"
)

func main() {
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))
	router := web.MainRouter()
	if strings.HasPrefix(Settings.HttpListen, "localhost") {
		webfatso.AppDemo(router)
	}

	// router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	if Settings.ResUrl == "/static/" {
		webfatso.ServStatic(router)
	}

	fmt.Printf("Start fat service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)
	err := http.ListenAndServe(Settings.HttpListen, router) // Start the server!
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
