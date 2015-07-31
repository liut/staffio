package main

import (
	"fmt"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	. "tuluu.com/liut/staffio/settings"
	_ "tuluu.com/liut/staffio/staffio-tight/statik"
	"tuluu.com/liut/staffio/web"
)

func main() {
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
