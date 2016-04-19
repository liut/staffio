package webfatso

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	_ "lcgc/platform/staffio/webfatso/statik"
)

func ServStatic(router *mux.Router) {

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
