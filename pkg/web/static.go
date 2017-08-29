package web

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	statikFs "github.com/rakyll/statik/fs"

	_ "github.com/liut/staffio/pkg/web/statik"
)

type assetsImpl struct {
	fs http.FileSystem
}

func newAssets(root, name string) *assetsImpl {
	fs := buildStaticFS(root, name)
	return &assetsImpl{fs}
}

func (a *assetsImpl) stripRouter(r gin.IRouter) {
	s := http.FileServer(a.fs)
	h := func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	}

	r.GET("/static/*filepath", h)
	r.GET("/favicon.ico", h)
	r.GET("/robots.txt", h)
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
