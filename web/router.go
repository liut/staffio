package web

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"tuluu.com/liut/staffio/backends"
	. "tuluu.com/liut/staffio/settings"
)

var (
	store              sessions.Store
	router             *mux.Router
	resUrl             string
	jsonRequestHeaders = []string{
		// "Accept", "application/json",
		"X-Requested-With", "XMLHttpRequest",
	}
	server *osin.Server
)

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rw, "Not Found")
}

func NewServerConfig() *osin.ServerConfig {
	return &osin.ServerConfig{
		AuthorizationExpiration: 900,
		AccessExpiration:        3600,
		TokenType:               "bearer",
		AllowedAuthorizeTypes: osin.AllowedAuthorizeType{
			osin.CODE,
			osin.TOKEN,
		},
		AllowedAccessTypes: osin.AllowedAccessType{
			osin.AUTHORIZATION_CODE,
			osin.IMPLICIT,
			osin.REFRESH_TOKEN,
			osin.PASSWORD,
			osin.CLIENT_CREDENTIALS,
		},
		ErrorStatusCode:           200,
		AllowClientSecretInParams: true,
		AllowGetAccessRequest:     false,
	}
}

func MainRouter() *mux.Router {
	if router != nil {
		return router
	}

	log.SetFlags(log.Ltime | log.Lshortfile)
	Settings.Parse()
	resUrl = Settings.ResUrl
	backends.Prepare()
	server = osin.NewServer(NewServerConfig(), backends.NewStorage())

	store = sessions.NewCookieStore([]byte(Settings.Session.Name))

	router = mux.NewRouter()

	router.Handle("/login", handler(loginForm)).Methods("GET").Name("login")
	router.Handle("/login", handler(login)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/logout", handler(logout)).Name("logout")
	router.Handle("/password", handler(passwordForm)).Methods("GET").Name("password")
	router.Handle("/password", handler(passwordChange)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/contacts", handler(contactListHandler)).Methods("GET")

	router.Handle("/authorize", handler(oauthAuthorize)).Methods("GET", "POST").Name("authorize")
	router.Handle("/token", handler(oauthToken)).Methods("GET", "POST").Name("token")
	router.Handle("/info", handler(oauthInfo)).Methods("GET", "POST").Name("info")

	router.Handle("/", handler(welcome)).Name("welcome")

	return router
}
