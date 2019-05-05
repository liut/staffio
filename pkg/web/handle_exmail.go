package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openshift/osin"
	"github.com/wealthworks/go-tencent-api/exmail"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
)

func (s *server) countNewMail(c *gin.Context) {
	user := UserWithContext(c)
	// log.Printf("user %q", user.UID)
	email := backends.GetEmailAddress(user.UID)
	res := make(osin.ResponseData)
	res["email"] = email

	count, err := exmail.CountNewMail(email)
	if err != nil {
		log.Printf("check new mail failed: %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	res["unseen"] = count
	res["got"] = true

	c.JSON(http.StatusOK, res)
}

func (s *server) loginToExmail(c *gin.Context) {
	user := UserWithContext(c)
	email := user.UID + "@" + settings.EmailDomain
	url, err := exmail.GetLoginURL(email)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}
