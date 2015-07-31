package web

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/goods/httpbuf"
	"log"
	"net/http"
	"net/url"
	"tuluu.com/liut/staffio/backends"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

// Authorization code endpoint
func oauthAuthorize(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {

		link := fmt.Sprintf("/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s",
			ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))
		// HANDLE LOGIN PAGE HERE
		if ctx.User == nil {
			ctx.Referer = link
			return loginForm(w, r, ctx)
			// resp.SetRedirect(reverse("login") + "?referer=" + reverse("authorize"))
		} else {
			if r.Method == "GET" {
				scopes, err := backends.LoadScopes()
				if err != nil {
					return err
				}
				return T("authorize.html").Execute(w, map[string]interface{}{
					"link":          link,
					"response_type": ar.Type,
					"scopes":        scopes,
					"client":        ar.Client.(*models.Client),
					"ctx":           ctx,
				})
			}

			if r.PostForm.Get("authorize") == "1" {
				ar.UserData = ctx.User.Uid
				ar.Authorized = true
				server.FinishAuthorizeRequest(resp, r, ar)
			} else {
				resp.SetRedirect(reverse("welcome"))
			}

		}

	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("authorize ERROR: %s\n", resp.InternalError)
	}
	// if !resp.IsError {
	// 	resp.Output["uid"] = ctx.User.Uid
	// }

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

// Access token endpoint
func oauthToken(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	var (
		uid  string = ""
		user *User
	)
	if ar := server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			uid = ar.UserData.(string)
			staff, err := backends.GetStaff(uid)
			if err != nil {
				resp.SetError("get_user_error", "staff not found")
				resp.InternalError = err
			} else {
				user = UserFromStaff(staff)
			}
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if Settings.HttpListen == "localhost:3000" && ar.Username == "test" && ar.Password == "test" {
				ar.Authorized = true
				break
			}

			if !backends.Authenticate(ar.Username, ar.Password) {
				resp.SetError("authentication_failed", err.Error())
				break
			}
			staff, err := backends.GetStaff(ar.Username)
			if err != nil {
				// resp.InternalError = err
				resp.SetError("get_user_failed", err.Error())
				break
			}
			ar.Authorized = true
			ar.UserData = staff.Uid
			user = UserFromStaff(staff)

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
		if uid != "" {
			resp.Output["uid"] = uid
		}
		if user != nil {
			resp.Output["user"] = user
		}

	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

// Information endpoint
func oauthInfo(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	var (
		uid  string
		user *User
	)
	if ir := server.HandleInfoRequest(resp, r); ir != nil {
		log.Printf("ir Code %s Token %s", ir.Code, ir.AccessData.AccessToken)
		uid = ir.AccessData.UserData.(string)
		staff, err := backends.GetStaff(uid)
		if err != nil {
			resp.SetError("get_user_error", "staff not found")
			resp.InternalError = err
		} else {
			user = UserFromStaff(staff)
		}
		server.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("info ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		if user != nil {
			resp.Output["user"] = user
		}
		resp.Output["uid"] = uid
	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

func welcome(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {

	log.Printf("session Name: %s, Values: %v", ctx.Session.Name(), ctx.Session.Values)
	log.Printf("ctx User %v", ctx.User)

	//execute the template
	return T("welcome.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func contactListHandler(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
	}
	limit := 5
	staffs := backends.ListPaged(limit)
	models.ByUid.Sort(staffs)

	return T("contact.html").Execute(w, map[string]interface{}{
		"staffs": staffs,
		"ctx":    ctx,
	})
}

func loginForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	return T("login.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func login(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	uid, password := req.FormValue("username"), req.FormValue("password")
	log.Printf("accept: %v (%d)", req.Header["Accept"], len(req.Header["Accept"]))
	res := make(osin.ResponseData)
	if !backends.Authenticate(uid, password) {
		// ctx.Session.AddFlash("Invalid Username/Password")

		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "password"}
		return outputJson(res, w)
	}

	staff, err := backends.GetStaff(uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Load user failed"}
		return outputJson(res, w)
	}

	//store the user id in the values and redirect to welcome
	ctx.Session.Values["user"] = UserFromStaff(staff)
	ctx.Session.Values["last_uid"] = staff.Uid

	res["ok"] = true
	res["referer"] = ctx.Referer
	return outputJson(res, w)
	// http.Redirect(w, req, reverse("welcome"), http.StatusSeeOther)
}

func logout(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	delete(ctx.Session.Values, "user")
	http.Redirect(w, req, reverse("welcome"), http.StatusSeeOther)
	return nil
}

func passwordForm(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	return T("password.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func passwordChange(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	uid, pwdOld, pwdNew := req.FormValue("username"), req.FormValue("old_password"), req.FormValue("new_password")
	res := make(osin.ResponseData)
	if !backends.Authenticate(uid, pwdOld) {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "old_password"}
		return outputJson(res, w)
	}
	err := backends.PasswordChange(uid, pwdOld, pwdNew)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "old_password"}
	} else {
		res["ok"] = true
	}

	return outputJson(res, w)
}

func outputJson(res map[string]interface{}, w http.ResponseWriter) error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Printf("json encoding error: %s", err)
	}
	return err
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
