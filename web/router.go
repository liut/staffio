package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/wealthworks/csmtp"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
)

var (
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

	csmtp.Host = Settings.SMTP.Host
	csmtp.Port = Settings.SMTP.Port
	csmtp.Name = Settings.SMTP.SenderName
	csmtp.From = Settings.SMTP.SenderEmail
	csmtp.Auth(Settings.SMTP.SenderPassword)

	if Settings.SentryDSN != "" {
		raven.SetDSN(Settings.SentryDSN)
	}

	resUrl = Settings.ResUrl
	backends.Prepare()
	server = osin.NewServer(NewServerConfig(), backends.NewStorage())
	var err error
	server.AccessTokenGen, err = getTokenGenJWT()
	if err != nil {
		panic(err)
	}

	sessionInit()
	router = mux.NewRouter()

	router.Handle("/login", handler(loginForm)).Methods("GET").Name("login")
	router.Handle("/login", handler(login)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/logout", handler(logout)).Name("logout")
	router.Handle("/password", handler(passwordForm)).Methods("GET").Name("password")
	router.Handle("/password", handler(passwordChange)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/profile", handler(profileForm)).Methods("GET").Name("profile")
	router.Handle("/profile", handler(profilePost)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/contacts", handler(contactsTable)).Methods("GET").Name("contacts")
	router.Handle("/staff/{uid:[a-z]+}", handler(staffForm)).Methods("GET").Name("staff")
	router.Handle("/staff/{uid:[a-z]+}", handler(staffPost)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/staff/{uid:[a-z]+}", handler(staffDelete)).Methods("DELETE").Headers(jsonRequestHeaders...)

	router.Handle("/authorize", handler(oauthAuthorize)).Methods("GET", "POST").Name("authorize")
	router.Handle("/token", handler(oauthToken)).Methods("GET", "POST").Name("token")
	router.Handle("/info/{topic}", handler(oauthInfo)).Methods("GET", "POST", "OPTIONS").Name("info")

	router.Handle("/dust/clients", handler(clientsForm)).Methods("GET").Name("clients")
	router.Handle("/dust/clients", handler(clientsPost)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/dust/scopes", handler(scopesForm)).Methods("GET", "POST").Name("scopes")
	router.Handle("/dust/_status/{topic:[a-z]+}{ext:(.json|.html|)}", handler(handleStatus)).Methods("GET").Name("status")

	router.Handle("/", handler(welcome)).Name("welcome")

	return router
}
