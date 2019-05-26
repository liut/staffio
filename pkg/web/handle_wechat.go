package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/models/random"
	"github.com/liut/staffio/pkg/settings"
)

const (
	stateLength = 31
	wxPrefix    = "https://open.work.weixin.qq.com"
	cKeyStateWX = "wxState"
)

func (s *server) wechatOAuth2Start(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if len(origin) == 0 {
		origin = settings.Current.BaseURL
	}
	callback := c.Request.FormValue("callback")
	if len(callback) == 0 {
		callback = "api/auth/wechat/callback"
	}
	state := random.GenString(stateLength)
	var (
		wxuri string
		inApp bool
		qs    = fmt.Sprintf("appid=%s&agentid=%d&redirect_uri=%s/%s&state=%s",
			s.wxAuth.CorpID(), settings.Current.WechatPortalAgentID, origin, callback, state)
	)

	ua := c.Request.UserAgent()
	if strings.Contains(ua, "wxwork/") { //  'wxwork/' | 'MicroMessenger/'
		// log.Printf("ua %q", ua)
		inApp = true
		wxuri = fmt.Sprintf("%s/connect/oauth2/authorize?%s", wxPrefix, qs)
	} else {
		wxuri = fmt.Sprintf("%s/wwopen/sso/qrConnect?%s", wxPrefix, qs)
	}
	sess := ginSession(c)
	sess.Set(cKeyStateWX, state)
	SessionSave(sess, c.Writer)
	apiOk(c, gin.H{
		"inapp": inApp,
		"wxuri": wxuri,
	}, 0)
}

func (s *server) wechatOAuth2Callback(c *gin.Context) {

	// if appid, corpId := c.Request.FormValue("appid"), s.wxAuth.CorpID(); appid != corpId {
	// 	log.Printf("incorrect appid %s=%s", appid, corpId)
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }
	code := c.Request.FormValue("code")
	if len(code) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	state := c.Request.FormValue("state")
	if len(state) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sess := ginSession(c)
	vState := sess.Get(cKeyStateWX)
	if vState == nil {
		log.Printf("state %s is expired", state)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if s := vState.(string); s != state {
		log.Printf("mismatch state %s=%s", state, s)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ou, err := s.wxAuth.GetOAuth2User(settings.Current.WechatPortalAgentID, code)
	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}
	log.Printf("auth2 with wechat work OK %v", ou)

	staff, err := s.service.Get(strings.ToLower(ou.UserID))
	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}
	signinStaffGin(c, staff)
	// OK
	if c.Request.Method == "POST" {
		apiOk(c, nil, 0)
	} else {
		c.Redirect(302, "/")
	}
}
