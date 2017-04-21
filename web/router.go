package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/coocood/freecache"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/wealthworks/csmtp"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
)

var (
	resUrl             string
	jsonRequestHeaders = []string{
		// "Accept", "application/json",
		"X-Requested-With", "XMLHttpRequest",
	}
	ws    *webImpl
	cache *freecache.Cache
)

type webImpl struct {
	*mux.Router
	osvr *osin.Server
	fs   http.FileSystem
}

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rw, "Not Found")
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, http.StatusText(http.StatusBadRequest))
}

func NewServerConfig() *osin.ServerConfig {
	return &osin.ServerConfig{
		AuthorizationExpiration: 900,
		AccessExpiration:        3600 * 24,
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

func New() *webImpl {
	ws = &webImpl{}

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
	ws.osvr = osin.NewServer(NewServerConfig(), backends.NewStorage())
	var err error
	ws.osvr.AccessTokenGen, err = getTokenGenJWT()
	if err != nil {
		panic(err)
	}

	cache = freecache.NewCache(Settings.CacheSize)
	sessionInit()
	router := mux.NewRouter()
	ws.Router = router

	router.Handle("/login", handler(loginForm)).Methods("GET").Name("login")
	router.Handle("/login", handler(login)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/logout", handler(logout)).Name("logout")
	router.Handle("/password/forgot", handler(passwordForgotForm)).Methods("GET").Name("password_forgot")
	router.Handle("/password/forgot", handler(passwordForgot)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/password/reset", handler(passwordResetForm)).Methods("GET").Name("password_reset")
	router.Handle("/password/reset", handler(passwordReset)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/password", handler(passwordForm)).Methods("GET").Name("password")
	router.Handle("/password", handler(passwordChange)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/profile", handler(profileForm)).Methods("GET").Name("profile")
	router.Handle("/profile", handler(profilePost)).Methods("POST").Headers(jsonRequestHeaders...)
	router.Handle("/email/unseen", handler(countNewMail)).Methods("GET").Name("unseen")
	router.Handle("/email/open", handler(loginToExmail)).Methods("GET").Name("email-open")

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
	router.Handle("/dust/status/{topic}", handler(handleStatus)).Methods("GET").Name("status")

	router.Handle("/article/{id}", handler(articleView)).Methods("GET").Name("article")
	router.Handle("/dust/articles", handler(articleForm)).Methods("GET").Name("article_form")
	router.Handle("/dust/articles", handler(articlePost)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/dust/links", handler(linksForm)).Methods("GET").Name("links_form")
	router.Handle("/dust/links", handler(linksPost)).Methods("POST").Headers(jsonRequestHeaders...)

	router.Handle("/cas/logout", handler(casLogout))
	router.Handle("/validate", handler(casValidateV1))
	router.Handle("/serviceValidate", handler(casValidateV2))

	router.Handle("/", handler(welcome)).Name("welcome")

	ws.ServStatic(Settings.Root, Settings.FS)
	return ws
}

func (ws *webImpl) Run(addr string) error {
	return http.ListenAndServe(addr, ws.Router)
}

func (ws *webImpl) HandleFunc(path string, f http.HandlerFunc) {
	ws.Router.HandleFunc(path, f)
}
