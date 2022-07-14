package htdocs

import (
	"embed"
	"net/http"
)

//go:embed favicon.ico robots.txt
//go:embed static/*
var uifs embed.FS

func Load(name string) (string, error) {
	if data, err := uifs.ReadFile(name); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func FS() embed.FS {
	return uifs
}

func Handler() http.Handler {
	return http.FileServer(http.FS(uifs))
}
