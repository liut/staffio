package web

import (
	"net/http"
	// "path/filepath"

	"github.com/liut/staffio/pkg/web/static"
)

type assetsImpl struct {
	fs string
}

func newAssets(name string) *assetsImpl {
	return &assetsImpl{fs: name}
}

func (a *assetsImpl) GetHandler() http.Handler {
	logger().Infow("accets", "fs", a.fs)
	if a.fs == "local" {
		return http.FileServer(http.Dir("./htdocs"))
	}

	return static.Server
}
