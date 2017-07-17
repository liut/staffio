package web

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/coocood/freecache"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/tj/go-debug"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/settings"
)

var (
	cache   *freecache.Cache
	keepers []string
	debug   = Debug("staffio:web")
)

type server struct {
	*gin.Engine
	service backends.Servicer
	osvr    *osin.Server
}

func New() *server {
	service := backends.NewService()
	osvr := osin.NewServer(newOsinConfig(), service.OSIN())
	var err error
	osvr.AccessTokenGen, err = getTokenGenJWT()
	if err != nil {
		panic(err)
	}

	svr := &server{
		Engine:  gin.New(),
		service: service,
		osvr:    osvr,
	}

	fmt.Println("gin in mode: ", gin.Mode())
	if settings.Debug {
		svr.Use(gin.Logger())
		svr.Use(gin.Recovery())
	} else {
		if settings.SentryDSN != "" {
			raven.SetDSN(settings.SentryDSN)
			onlyCrashes := false
			svr.Use(sentry.Recovery(raven.DefaultClient, onlyCrashes))
		}
	}
	store := sessionStore()
	svr.Use(sessions.Sessions("session", store))
	svr.strapRouter(svr.Engine)

	cache = freecache.NewCache(settings.CacheSize)
	group, err := svr.service.GetGroup("keeper")
	if err != nil {
		panic(err)
	}
	keepers = group.Members

	return svr
}

func (s *server) HandleFunc(path string, hf http.HandlerFunc) {
	h := func(c *gin.Context) {
		hf(c.Writer, c.Request)
	}
	s.GET(path, h)
	s.POST(path, h)
}

func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s.Engine)
}

func newOsinConfig() *osin.ServerConfig {
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
