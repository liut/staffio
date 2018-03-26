package web

import (
	"log"
	"net/http"
	"path/filepath"

	statikFs "github.com/rakyll/statik/fs"

	_ "github.com/liut/staffio/pkg/web/statik"
)

type assetsImpl struct {
	fs   http.FileSystem
	Base string // url prefix
}

func newAssets(root, name string) *assetsImpl {
	fs := buildStaticFS(root, name)
	return &assetsImpl{fs: fs}
}

func (a *assetsImpl) GetHandler() http.Handler {
	h := http.FileServer(a.fs)
	if a.Base != "" && a.Base != "/" {
		return http.StripPrefix(a.Base, h)
	}
	return h
}

func buildStaticFS(root, name string) (fs http.FileSystem) {

	var (
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

	return
}
