package web

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"tuluu.com/liut/staffio/backends"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

type User struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

func (u *User) IsKeeper() bool {
	if u == nil {
		return false
	}
	keeper := backends.GetGroup("keeper")
	return keeper.Has(u.Uid)
}

func UserFromStaff(staff *models.Staff) *User {
	return &User{staff.Uid, staff.Name()}
}

type Context struct {
	Session   *sessions.Session
	ResUrl    string
	User      *User
	LastUid   string
	NavSimple bool
	Referer   string
}

func (c *Context) Close() {
	backends.CloseAll()
}

func NewContext(req *http.Request) (*Context, error) {
	sess, err := store.Get(req, Settings.Session.Name)
	sess.Options.Domain = Settings.Session.Domain
	sess.Options.HttpOnly = true
	var (
		lastUid string
		user    *User
	)
	if v, ok := sess.Values["last_uid"]; ok {
		lastUid = v.(string)
	}
	if v, ok := sess.Values["user"]; ok {
		user = v.(*User)
	}
	referer := req.FormValue("referer")
	if referer == "" {
		referer = req.Referer()
	}
	ctx := &Context{
		Session: sess,
		ResUrl:  Settings.ResUrl,
		Referer: referer,
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
