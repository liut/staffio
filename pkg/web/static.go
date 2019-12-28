package web

import (
	"net/http"
	// "path/filepath"

	"github.com/liut/staffio/pkg/web/static"
)

func staticHandler(fs string) http.Handler {
	if fs == "local" {
		return http.FileServer(http.Dir("./htdocs"))
	}

	return static.Server
}
