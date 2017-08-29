package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/settings"
)

func (s *server) strapRouter(r gin.IRouter) {

	r.GET("/login", s.loginForm).POST("/login", s.loginPost)
	r.GET("/logout", s.logout)
	r.GET("/password/forgot", s.passwordForgotForm)
	r.POST("/password/forgot", s.passwordForgot)
	r.GET("/password/reset", s.passwordResetForm)
	r.POST("/password/reset", s.passwordReset)

	authed := r.Group("/", AuthUserMiddleware())
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
	r.GET("/token", s.oauth2Token)
	r.POST("/token", s.oauth2Token)
	r.GET("/info/:topic", s.oauth2Info)
	r.POST("/info/:topic", s.oauth2Info)

	keeper := authed.Group("/dust", AuthAdminMiddleware())
	keeper.GET("/clients", s.clientsForm)
	keeper.POST("/clients", s.clientsPost)
	keeper.GET("/scopes", s.scopesForm)
	keeper.GET("/status/:topic", s.handleStatus)
	keeper.GET("/groups", s.groupList)

	r.GET("/article/:id", articleView)
	keeper.GET("/articles", articleForm)
	keeper.POST("/articles", articlePost)

	keeper.GET("/links", linksForm)
	keeper.POST("/links", linksPost)

	r.GET("/cas/logout", casLogout)
	r.GET("/validate", s.casValidateV1)
	r.GET("/serviceValidate", s.casValidateV2)

	r.GET("/", welcome)

	assets := newAssets(settings.Root, settings.FS)
	assets.stripRouter(r)
}

func IsAjax(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Index(accept, "application/json") >= 0
}
