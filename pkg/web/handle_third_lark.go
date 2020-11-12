package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"fhyx.online/lark-api-go/lark"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/random"
	"github.com/liut/staffio/pkg/settings"
)

const (
	larkPrefix  = "https://open.feishu.cn"
	cKeyStateLk = "lkState"
)

var larkDecr = lark.NewCrypto(settings.Current.LarkEncryptKey)

func inAPPLark(req *http.Request) bool {
	return strings.Contains(req.UserAgent(), "Lark/")
}

func (s *server) larkOAuth2Start(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if len(origin) == 0 {
		origin = settings.Current.BaseURL
	}
	callback := c.Request.FormValue("callback")
	if len(callback) == 0 {
		callback = "api/auth/feishu/callback"
	}
	state := random.GenString(stateLength)

	var (
		uri   string
		inApp bool
	)
	inApp = inAPPLark(c.Request)
	// uri = fmt.Sprintf("%s/open-apis/authen/v1/index?app_id=%s&redirect_uri=%s/%s&state=%s", larkPrefix, settings.Current.LarkAppID, origin, callback, state)
	qs := fmt.Sprintf("app_id=%s&redirect_uri=%s/%s&state=%s", settings.Current.LarkAppID, origin, callback, state)
	if inApp {
		uri = fmt.Sprintf("%s/open-apis/authen/v1/index?%s", larkPrefix, qs)
	} else {
		uri = fmt.Sprintf("%s/connect/qrconnect/page/sso/?%s", larkPrefix, qs)
	}
	logger().Infow("auth from lark", "ua", c.Request.UserAgent(), "uri", uri)

	sess := ginSession(c)
	sess.Set(cKeyStateLk, state)
	SessionSave(sess, c.Writer)
	apiOk(c, gin.H{
		"inapp": inApp,
		"fsuri": uri,
	}, 0)
}

func (s *server) larkOAuth2Callback(c *gin.Context) {
	code := c.Request.FormValue("code")
	if len(code) == 0 {
		logger().Infow("empty code")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	state := c.Request.FormValue("state")
	if len(state) == 0 {
		logger().Infow("empty state")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sess := ginSession(c)
	vState := sess.Get(cKeyStateLk)
	if vState == nil {
		logger().Infow("state is expired", "state", state)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if s, err := cast.ToStringE(vState); err != nil || s != state {
		logger().Infow("mismatch state", "state", state, "str", s)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ou, err := s.larkAPI.AuthorizeCode(code)
	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}
	logger().Infow("auth2 with lark callback OK", "ou", ou)

	data := s.service.All(&models.Spec{Name: ou.Name})
	if len(data) == 0 {
		logger().Infow("get staff fail", "name", ou.Name)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logger().Infow("found ", "data", data)
	signinStaffGin(c, &data[0])
	// OK
	if c.Request.Method == "POST" {
		apiOk(c, nil, 0)
	} else {
		// c.Redirect(302, "/")
		c.Writer.Header().Set("refresh", "1;url=/")
		fmt.Fprintf(c.Writer, refreshHTML, ou.Name)
	}
}

const (
	refreshHTML = `<html>
	Welcome %s, now loading...
	<script>parent.location.href='/';
	</script>
	</html>`
)

func (s *server) larkEventCallback(c *gin.Context) {
	var data lark.EncryptEntry
	if err := c.Bind(&data); err != nil {
		logger().Infow("bind fail", "err", err)
		apiError(c, 400, err)
		return
	}
	// logger().Infow("decrypt ok", "body", data.EncryptedBody)

	decryptedText, err := larkDecr.DecryptString(data.EncryptedBody)
	if err != nil {
		logger().Infow("decrypt fail", "err", err)
		apiError(c, 400, err)
		return
	}

	cr, err := lark.UnmarsalCallback(decryptedText)
	if err != nil {
		logger().Infow("unmarshal fail", "decryptedText", decryptedText, "err", err)
		apiError(c, 400, err)
		return
	}

	if cr.Type == "url_verification" {
		c.JSON(200, lark.CallbackResp{Challenge: cr.Challenge})
		return
	}
	logger().Infow("got lark event callback", "type", cr.Type)

}
