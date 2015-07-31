package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	. "tuluu.com/liut/staffio/settings"
	"tuluu.com/liut/staffio/web"
	"tuluu.com/liut/staffio/webfatso"
)

func main() {
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
