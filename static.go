package main

import (
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	_ "tuluu.com/liut/staffio/statik"
)

func staticServ() {

	statikFS, se := fs.New()
	if se != nil {
		log.Fatalf(se.Error())
	}

	// statikFS := http.Dir(filepath.Join(Settings.Root, "htdocs"))
	ss := http.FileServer(statikFS)
	router.PathPrefix("/static/").Handler(ss).Methods("GET", "HEAD")
	router.Path("/favicon.ico").Handler(ss).Methods("GET", "HEAD")
	router.Path("/robots.txt").Handler(ss).Methods("GET", "HEAD")

}
