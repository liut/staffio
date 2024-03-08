package web

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/settings"
	xrefs "github.com/liut/staffio/pkg/xrefs"
)

var (
	avatarReplacer = strings.NewReplacer("/0", "/60") //nolint
)

const (
	kReferer = "_ref"
)

func refererWithContext(c *gin.Context) (ref string) {
	cookie, err := c.Cookie(kReferer)
	if err == nil {
		ref = cookie
	}
	return
}

func markReferer(c *gin.Context) {
	c.SetCookie(kReferer, c.Request.RequestURI, 10, "/", "", false, true)
}

func (s *server) Render(c *gin.Context, name string, data map[string]any) {
	instance := s.tpl(name)
	m := data
	m["base"] = base
	m["appVersion"] = settings.Version()
	m["navSimple"] = false
	session := ginSession(c)
	m["session"] = session.Values()
	m["referer"] = refererWithContext(c)
	var user *User
	v, exist := c.Get(kAuthUser)
	if exist {
		user = v.(*User)
	} else {
		user, _ = authzr.UserFromRequest(c.Request)
	}
	m["currUser"] = user
	m["checkEmail"] = settings.Current.EmailCheck
	if err := instance.Execute(c.Writer, m); err != nil {
		logger().Infow("render fail", "err", err)
	}
}

func (s *server) tpl(name string) *template.Template {

	t := template.New("_base.html").Funcs(template.FuncMap{
		"urlFor":     UrlFor,
		"avatarHtml": AvatarHTML,
		"isKeeper":   s.IsKeeper,
	})
	t = template.Must(t.ParseFS(
		xrefs.FS(),
		"templates/_base.html",
		"templates/"+name,
	))

	return t
}

// AvatarHTML 生成头像的HTML标签，目前仅支持微信头像
func AvatarHTML(s string) template.HTML {
	if len(s) == 0 {
		return ""
	}
	return template.HTML("<img class='avatar' src=\"" + s + "\">")
}
