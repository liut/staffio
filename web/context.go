package web

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
)

type Context struct {
	Request   *http.Request
	Vars      map[string]string
	Session   *sessions.Session
	ResUrl    string
	User      *User
	LastUid   string
	NavSimple bool
	Referer   string
	Version   string
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
)

func NewContext(req *http.Request) (*Context, error) {
	sess, err := store.Get(req, Settings.Session.Name)
	sess.Options.Domain = Settings.Session.Domain
	sess.Options.HttpOnly = true
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
		referer = req.Referer()
	}
	ctx := &Context{
		Request: req,
		Vars:    mux.Vars(req),
		Session: sess,
		ResUrl:  Settings.ResUrl,
		Referer: referer,
		Version: Settings.Version,
		LastUid: lastUid,
		User:    user,
	}
	if err != nil {
		log.Printf("new context error: %s", err)
		return ctx, err
	}

	return ctx, err
}

func init() {
	gob.Register(&User{})
}
