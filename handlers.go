package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	// "github.com/RangelReale/osin"
	"github.com/goods/httpbuf"
	"log"
	"net/http"
	"tuluu.com/liut/staffio/backends/ldap"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

var (
	backendReady bool
)

func prepareBackends() {
	if backendReady {
		return
	}
	addr := fmt.Sprintf("%s:%d", Settings.LDAP.Host, Settings.LDAP.Port)
	ls := ldap.AddSource(addr, Settings.LDAP.Base)
	ls.BindDN = Settings.LDAP.BindDN
	ls.Passwd = Settings.LDAP.Password

	backendReady = true
}

func contactListHandler(rw http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	prepareBackends()
	limit := 5
	staffs := ldap.ListPaged(limit)

	return T("contact.html").Execute(rw, map[string]interface{}{
		"staffs": staffs,
		"ctx":    ctx,
	})
}

type Context struct {
	Session   *sessions.Session
	ResUrl    string
	User      *models.Staff
	NavSimple bool
}

func (c *Context) Close() {
	// c.Session
}

func NewContext(req *http.Request) (*Context, error) {
	sess, err := store.Get(req, Settings.Session.Name)
	ctx := &Context{
		Session: sess,
		ResUrl:  Settings.ResUrl,
	}
	if err != nil {
		log.Printf("new context error: %s", err)
		return ctx, err
	}

	return ctx, err
}

type handler func(http.ResponseWriter, *http.Request, *Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := NewContext(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	//run the handler and grab the error, and report it
	buf := new(httpbuf.Buffer)
	err = h(buf, req, ctx)
	if err != nil {
		log.Printf("call handler error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//save the session
	if err = ctx.Session.Save(req, buf); err != nil {
		log.Printf("session.save error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//apply the buffered response to the writer
	buf.Apply(w)
}
