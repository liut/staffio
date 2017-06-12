package web

import (
	"log"
	"net/http"
	"path/filepath"

	statikFs "github.com/rakyll/statik/fs"

	_ "lcgc/platform/staffio/pkg/web/statik"
)

func (ws *webImpl) ServStatic(root, name string) {

	var (
		fs  http.FileSystem
		err error
	)
	log.Printf("using fs %s", name)
	if name != "local" {
		fs, err = statikFs.New()
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		fs = http.Dir(filepath.Join(root, "htdocs"))
	}

	h := http.FileServer(fs)
	ws.PathPrefix("/static/").Handler(h).Methods("GET", "HEAD")
	ws.Path("/favicon.ico").Handler(h).Methods("GET", "HEAD")
	ws.Path("/robots.txt").Handler(h).Methods("GET", "HEAD")

}
