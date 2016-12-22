package web

import (
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
)

type Context struct {
	Request   *http.Request
	Writer    http.ResponseWriter
	Vars      map[string]string
	Session   *sessions.Session
	ResUrl    string
	User      *User
	LastUid   string
	NavSimple bool
	Referer   string
	Version   string
}

func (c *Context) IsUserExpired() bool {
	return c.User == nil || c.User.IsExpired()
}

func (c *Context) afterHandle() {
	if c.User != nil {
		if !c.User.IsExpired() {
			c.User.Refresh()
		}
	}
}

func (c *Context) Close() {
	backends.CloseAll()
}

const (
	kLastUid = "lu"
	kUserOL  = "user"
	kRefer   = "ref"
)

func NewContext(w http.ResponseWriter, req *http.Request, sess *sessions.Session) *Context {
	var (
		lastUid string
		user    *User
	)
	if v, ok := sess.Values[kLastUid]; ok {
		lastUid = v.(string)
	}
	if v, ok := sess.Values[kUserOL]; ok {
		user = v.(*User)
	}
	referer := req.FormValue("referer")
	if referer == "" {
		if ref, ok := sess.Values[kRefer]; ok {
			referer = ref.(string)
		}
		if referer == "" {
			referer = req.Referer()
		}
	}
	// log.Printf("sessions %v", sess.Values)
	ctx := &Context{
		Request: req,
		Writer:  w,
		Vars:    mux.Vars(req),
		Session: sess,
		ResUrl:  Settings.ResUrl,
		Referer: referer,
		Version: Settings.Version,
		LastUid: lastUid,
		User:    user,
	}

	return ctx
}

func (ctx *Context) checkLogin() bool {
	if ctx.IsUserExpired() {
		ctx.toLogin()
		return false
	}
	return true
}

func (ctx *Context) toLogin() {
	ctx.Session.Values[kRefer] = ctx.Request.RequestURI
	http.Redirect(ctx.Writer, ctx.Request, reverse("login"), http.StatusTemporaryRedirect)
}

func (ctx *Context) Render(tpl string, data interface{}) error {
	return T(tpl).Execute(ctx.Writer, data)
}

func (ctx *Context) IsAjax() bool {
	accept := ctx.Request.Header.Get("Accept")
	return strings.Index(accept, "application/json") >= 0
}

func init() {
	gob.Register(&User{})
}
