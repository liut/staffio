package web

import (
	"fmt"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"
	"github.com/openshift/osin"

	"github.com/fhyx/lark-api-go/lark"
	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
)

var (
	svr *server
)

// Config ...
type Config struct {
	Addr    string
	Root    string
	FS      string
	BaseURI string
}

type server struct {
	cfg     Config
	router  *gin.Engine
	service backends.Servicer
	osvr    *osin.Server
	wxAuth  *exwechat.API
	checkin *exwechat.CAPI
	larkAPI *lark.API
}

func (s *server) IsKeeper(uid string) bool {
	return s.InGroup("keeper", uid)
}

func (s *server) InGroup(gn, uid string) bool {
	return s.service.InGroup(gn, uid)
}

func (s *server) InGroupAny(uid string, gn ...string) bool {
	return s.service.InGroupAny(uid, gn...)
}

// New returns current server instance
func New(c Config) *server {
	if svr != nil {
		return svr
	}
	service := backends.NewService()

	// check ready
	if err := service.Ready(); err != nil {
		panic(err)
	}

	osvr := osin.NewServer(newOsinConfig(), service.OSIN())
	var err error
	osvr.AccessTokenGen, err = getTokenGenJWT()
	if err != nil {
		panic(err)
	}

	svr = &server{
		cfg:     c,
		router:  gin.New(),
		service: service,
		osvr:    osvr,
		wxAuth:  exwechat.New(settings.Current.WechatCorpID, settings.Current.WechatPortalSecret),
		checkin: exwechat.NewCAPI(),
		larkAPI: lark.NewAPI(settings.Current.LarkAppID, settings.Current.LarkAppSecret),
	}

	if settings.IsDevelop() {
		fmt.Printf("In Developing(Debug) mode, gin: %s\n", gin.Mode())
		svr.router.Use(gin.Logger(), gin.Recovery())
	} else {
		fmt.Printf("In Release mode, gin: %s\n", gin.Mode())
		if settings.Current.SentryDSN != "" {
			raven.SetDSN(settings.Current.SentryDSN)
			onlyCrashes := false
			svr.router.Use(sentry.Recovery(raven.DefaultClient, onlyCrashes))
		}
	}

	svr.StrapRouter()

	return svr
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: refactory
	s.router.ServeHTTP(w, req)
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
			// osin.REFRESH_TOKEN,
			osin.PASSWORD,
			// osin.CLIENT_CREDENTIALS,
		},
		ErrorStatusCode:           200,
		AllowClientSecretInParams: true,
		AllowGetAccessRequest:     false,
	}
}
