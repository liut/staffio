package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/settings"
)

var (
	base = "/"
)

func SetBase(s string) {
	base = fmt.Sprintf("%s/", strings.TrimRight(s, "/"))
}

func (s *server) StrapRouter() {
	gr := s.router.Group(base)
	gr.GET("/login", s.loginForm).POST("/login", s.loginPost)
	gr.GET("/logout", s.logout)
	gr.GET("/password/forgot", s.passwordForgotForm)
	gr.POST("/password/forgot", s.passwordForgot)
	gr.GET("/password/reset", s.passwordResetForm)
	gr.POST("/password/reset", s.passwordReset)

	authed := gr.Group("/", AuthUserMiddleware())
	authed.GET("/password", s.passwordForm)
	authed.POST("/password", s.passwordChange)

	authed.GET("/profile", s.profileForm)
	authed.POST("/profile", s.profilePost)
	authed.GET("/email/unseen", s.countNewMail)
	authed.GET("/email/open", s.loginToExmail)

	authed.GET("/contacts", s.contactsTable)
	authed.GET("/staff/:uid", s.staffForm)
	authed.POST("/staff/:uid", s.staffPost)
	authed.DELETE("/staff/:uid", s.staffDelete)

	authed.GET("/authorize", s.oauth2Authorize)
	authed.POST("/authorize", s.oauth2Authorize)
	gr.GET("/token", s.oauth2Token)
	gr.POST("/token", s.oauth2Token)
	gr.GET("/info/:topic", s.oauth2Info)
	gr.POST("/info/:topic", s.oauth2Info)

	keeper := authed.Group("/dust", AuthAdminMiddleware())
	keeper.GET("/clients", s.clientsForm)
	keeper.POST("/clients", s.clientsPost)
	keeper.GET("/scopes", s.scopesForm)
	keeper.GET("/status/:topic", s.handleStatus)
	keeper.GET("/groups", s.groupList)

	gr.GET("/article/:id", articleView)
	keeper.GET("/articles", articleForm)
	keeper.POST("/articles", articlePost)

	keeper.GET("/links", linksForm)
	keeper.POST("/links", linksPost)

	gr.GET("/cas/logout", casLogout)
	gr.GET("/validate", s.casValidateV1)
	gr.GET("/serviceValidate", s.casValidateV2)

	gr.GET("/", welcome)

	assets := newAssets(settings.Root, settings.FS)
	assets.Base = base
	ah := gin.WrapH(assets.GetHandler())
	s.router.GET("/static/*filepath", ah)
	s.router.GET("/favicon.ico", ah)
	s.router.GET("/robots.txt", ah)

}

func IsAjax(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Index(accept, "application/json") >= 0
}

func UrlFor(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(path, "/"))
}
