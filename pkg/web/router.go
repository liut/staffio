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

	authed := gr.Group("/", AuthUserMiddleware(true))
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

	keeper := authed.Group("/dust", s.authAdminMiddleware())
	keeper.GET("/clients", s.clientsForm)
	keeper.POST("/clients", s.clientsPost)
	keeper.GET("/scopes", s.scopesForm)
	keeper.GET("/status/:topic", s.handleStatus)
	keeper.GET("/groups", s.groupList)

	gr.GET("/article/:id", s.articleView)
	keeper.GET("/articles", s.articleForm)
	keeper.POST("/articles", s.articlePost)

	keeper.GET("/links", s.linksForm)
	keeper.POST("/links", s.linksPost)

	gr.GET("/cas/logout", casLogout)
	gr.GET("/validate", s.casValidateV1)
	gr.GET("/serviceValidate", s.casValidateV2)

	gr.GET("/", s.welcome)

	{ // for new lcgc/staff only
		gr.GET("/api/me", s.me)
		gr.POST("/api/verify", s.me)
		gr.POST("/api/login", s.loginPost)
		gr.POST("/api/logout", s.logout)
		gr.POST("/api/password/forgot", s.passwordForgot)
		gr.POST("/api/password/reset", s.passwordReset)
	}

	api := gr.Group("/api", AuthUserMiddleware(false))
	{
		api.POST("/weekly/report/add", s.weeklyReportAdd)
		api.POST("/weekly/report/update", s.weeklyReportUpdate)
		api.POST("/weekly/report/up", s.weeklyReportUp)
		api.POST("/weekly/report/all", s.weeklyReportList)
		api.POST("/weekly/report/self", s.weeklyReportListSelf)
		api.POST("/weekly/problems", s.weeklyProblemList)
		api.POST("/weekly/problem/add", s.weeklyProblemAdd)
		api.POST("/weekly/problem/update", s.weeklyProblemUpdate)
		api.GET("/staffs", s.staffList)
		api.GET("/teams", s.teamListByRole)
		api.POST("/team/member", s.teamMemberOp)

		apiMan := api.Group("/", s.authAdminMiddleware())
		apiMan.POST("/weekly/report/stat", s.weeklyReportStat)
		apiMan.POST("/weekly/report/ignore/add", s.weeklyIgnoreAdd)
		apiMan.POST("/weekly/report/ignore/del", s.weeklyIgnoreRemove)
		apiMan.GET("/weekly/report/ignores", s.weeklyIgnoreList)
		apiMan.GET("/weekly/report/vacations", s.weeklyVacationList)
		apiMan.POST("/weekly/report/vacation/mark", s.weeklyVacationAdd)
		apiMan.POST("/weekly/report/vacation/unmark", s.weeklyVacationRemove)
	}

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

func apiError(c *gin.Context, status int, message interface{}) {
	resp := map[string]interface{}{
		"status": status,
	}
	switch ret := message.(type) {
	case error:
		resp["message"] = ret.Error()
	default:
		resp["message"] = ret
	}
	c.JSON(http.StatusOK, resp)
}

func apiOk(c *gin.Context, data interface{}, count int) {
	res := map[string]interface{}{"status": 0}
	if data != nil {
		res["data"] = data
	}
	if count > 0 {
		res["count"] = count
	}
	c.JSON(http.StatusOK, res)
}
