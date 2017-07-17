package web

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin"
	"github.com/wealthworks/go-tencent-api/exmail"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/settings"
)

func (s *server) countNewMail(c *gin.Context) {
	user := UserWithContext(c)
	// log.Printf("user %q", user.Uid)
	email := backends.GetEmailAddress(user.Uid)
	res := make(osin.ResponseData)
	res["email"] = email
	key := []byte(fmt.Sprintf("mail-count-%s", user.Uid))

	if bv, err := cache.Get(key); err == nil {
		res["unseen"] = binary.LittleEndian.Uint32(bv)
	} else {
		count, err := exmail.CountNewMail(email)
		if err != nil {
			log.Printf("check new mail failed: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(count))
		cache.Set(key, bs, int(settings.CacheLifetime))
		res["unseen"] = count
		res["got"] = true
	}

	c.JSON(http.StatusOK, res)
}

func (s *server) loginToExmail(c *gin.Context) {
	user := UserWithContext(c)
	email := user.Uid + "@" + settings.EmailDomain
	url, err := exmail.GetLoginURL(email)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}
