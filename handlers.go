package main

import (
	"fmt"
	// "github.com/gorilla/sessions"
	"github.com/RangelReale/osin"
	"github.com/goods/httpbuf"
	"log"
	"net/http"
	"tuluu.com/liut/staffio/backends"
	// "tuluu.com/liut/staffio/models"
	// . "tuluu.com/liut/staffio/settings"
	"encoding/json"
	"net/url"
)

// Authorization code endpoint
func oauthAuthorize(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {

		// HANDLE LOGIN PAGE HERE
		if ctx.User == nil {
			referer := fmt.Sprintf("/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s",
				ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))
			ctx.Referer = referer
			return loginForm(w, r, ctx)
			// resp.SetRedirect(reverse("login") + "?referer=" + reverse("authorize"))
		} else {
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}

	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("authorize ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 19923
	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

// Access token endpoint
func oauthToken(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if ar.Username == "test" && ar.Password == "test" {
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		case osin.ASSERTION:
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}
		server.FinishAccessRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("token ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 19923
	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

// Information endpoint
func oauthInfo(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ir := server.HandleInfoRequest(resp, r); ir != nil {
		server.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("info ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 19923
	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

func index(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {

	log.Printf("session Name: %s, Values: %v", ctx.Session.Name(), ctx.Session.Values)
	log.Printf("ctx User %v", ctx.User)

	//execute the template
	return T("index.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func contactListHandler(rw http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	backends.Prepare()
	limit := 5
	staffs := backends.ListPaged(limit)
	keeper := backends.GetGroup("keeper")

	return T("contact.html").Execute(rw, map[string]interface{}{
		"staffs": staffs,
		"keeper": keeper,
		"ctx":    ctx,
	})
}

func loginForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	return T("login.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func login(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	username, password := req.FormValue("username"), req.FormValue("password")

	user, e := backends.Login(username, password)
	if e != nil {
		// ctx.Session.AddFlash("Invalid Username/Password")

		res := make(osin.ResponseData)
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "password"}
		outputJson(res, w)
		return nil
	}

	//store the user id in the values and redirect to index
	ctx.Session.Values["user"] = &User{user.Uid, user.Name()}

	res := make(osin.ResponseData)
	res["ok"] = true
	res["referer"] = reverse("index")
	outputJson(res, w)
	// http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return nil
}

func logout(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	delete(ctx.Session.Values, "user")
	http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return nil
}

func outputJson(res map[string]interface{}, w http.ResponseWriter) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Printf("json encoding error: %s", err)
	}
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
