package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	. "tuluu.com/liut/staffio/settings"
)

var (
	store  sessions.Store
	router *mux.Router
)

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rw, "Not Found")
}

func main() {

	Settings.Parse()
	prepareBackends()

	store = sessions.NewCookieStore([]byte(Settings.Session.Name))

	router = mux.NewRouter()
	router.Handle("/contact", handler(contactListHandler)).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	fmt.Printf("Start service %s at addr %s\n", Settings.Version, Settings.HttpListen)
	err := http.ListenAndServe(Settings.HttpListen, router) // Start the server!
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
