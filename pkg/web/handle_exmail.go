package web

import (
	"net/http"

	"daxv.cn/gopak/tencent-api-go/exmail"
	"github.com/gin-gonic/gin"
	"github.com/go-osin/osin"

	"github.com/liut/staffio/pkg/backends/qqexmail"
	"github.com/liut/staffio/pkg/settings"
)

// TODO: upgrade to new api of wework
func (s *server) countNewMail(c *gin.Context) {
	user := UserWithContext(c)
	email := qqexmail.GetEmailAddress(user.UID)
	res := make(osin.ResponseData)
	res["email"] = email

	count, err := exmail.CountNewMail(email)
	if err != nil {
		logger().Infow("check new mail fail", "err", err)
		c.AbortWithError(http.StatusInternalServerError, err) //nolint
		return
	}
	res["unseen"] = count
	res["got"] = true

	c.JSON(http.StatusOK, res)
}

func (s *server) loginToExmail(c *gin.Context) {
	user := UserWithContext(c)
	email := user.UID + "@" + settings.Current.EmailDomain
	url, err := exmail.GetLoginURL(email)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err) //nolint
		return
	}
	c.Redirect(http.StatusFound, url)
}
