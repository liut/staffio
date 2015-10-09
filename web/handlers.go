package web

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/goods/httpbuf"
	"log"
	"net/http"
	"net/url"
	"strings"
	"tuluu.com/liut/staffio/backends"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

// Authorization code endpoint
func oauthAuthorize(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
		link := fmt.Sprintf("/authorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
			ar.Type, ar.Client.GetId(), url.QueryEscape(ar.RedirectUri), ar.State, ar.Scope)
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

	debugf("oauthAuthorize resp: %v", resp)
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
		debugf("ar Code %s Scope %s", ar.Code, ar.Scope)
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
				ar.UserData = "test"
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
			resp.Output["is_keeper"] = IsKeeper(uid)
		}
		if user != nil {
			resp.Output["user"] = user
		}

	}

	debugf("oauthToken resp: %v", resp)

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

// Information endpoint
func oauthInfo(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	if ir := server.HandleInfoRequest(resp, r); ir != nil {
		debugf("ir Code %s Token %s", ir.Code, ir.AccessData.AccessToken)
		var (
			uid   string
			topic = ctx.Vars["topic"]
		)
		uid = ir.AccessData.UserData.(string)
		staff, err := backends.GetStaff(uid)
		if err != nil {
			resp.SetError("get_user_error", "staff not found")
			resp.InternalError = err
		} else {
			resp.Output["uid"] = uid
			if strings.HasPrefix(topic, "me") {
				resp.Output["me"] = staff
				if len(topic) > 2 && strings.Index(topic, "+") == 2 {
					// TODO: search group topic[2:]
				}
			} else if topic == "staff" {
				resp.Output["staff"] = staff
			}

		}
		server.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("info ERROR: %s\n", resp.InternalError)
	}

	osin.OutputJSON(resp, w, r)
	return resp.InternalError
}

func clientsForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil || !ctx.User.IsKeeper() {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}
	var (
		limit  = 20
		offset = 0
		sort   = map[string]int{"id": backends.ASCENDING}
	)
	clients, err := backends.LoadClients(limit, offset, sort)
	if err != nil {
		return err
	}
	return T("clients.html").Execute(w, map[string]interface{}{
		"ctx":     ctx,
		"clients": clients,
	})
}

func clientsPost(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil || !ctx.User.IsKeeper() {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}
	res := make(osin.ResponseData)

	var (
		client *models.Client
	)

	if req.FormValue("op") == "new" {
		// create new client
		client = models.NewClient(
			req.PostFormValue("name"),
			req.PostFormValue("code"),
			req.PostFormValue("secret"),
			req.PostFormValue("redirect_uri"))
		// log.Printf("new client: %v", client)
		_, e := backends.GetClientWithCode(client.Code) // check exists
		if e == nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": "duplicate client_id"}
			return outputJson(res, w)
		}

	} else {

		pk, name, value := req.PostFormValue("pk"), req.PostFormValue("name"), req.PostFormValue("value")
		log.Printf("clientsPost: pk %s, name %s, value %s", pk, name, value)
		if pk == "" {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is empty"}
			return outputJson(res, w)
		}
		// id, err := strconv.ParseUint(pk, 10, 32)
		client, err = backends.GetClientWithCode(pk)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is invalid or not found"}
			return outputJson(res, w)
		}
		switch name {
		case "name":
			client.Name = value
		case "secret":
			client.Secret = value
		case "redirect_uri":
			client.RedirectUri = value
		default:
			log.Printf("invalid filed: %s", name)
		}
	}

	if client != nil {
		err = backends.SaveClient(client)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
			return outputJson(res, w)
		}
		res["ok"] = true
		res["id"] = client.Id
		return outputJson(res, w)
	}

	res["ok"] = false
	res["error"] = map[string]string{"message": "invalid operation"}
	return outputJson(res, w)
}

func scopesForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil || !ctx.User.IsKeeper() {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}
	scopes, err := backends.LoadScopes()
	if err != nil {
		return err
	}
	return T("scopes.html").Execute(w, map[string]interface{}{
		"ctx":    ctx,
		"scopes": scopes,
	})
	return nil
}

func welcome(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {

	if Settings.Debug {
		log.Printf("session Name: %s, Values: %d", ctx.Session.Name(), len(ctx.Session.Values))
		log.Printf("ctx User %v", ctx.User)
	}

	//execute the template
	return T("welcome.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
}

func contactsTable(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}

	staffs := backends.LoadStaffs()
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
	// log.Printf("accept: %v (%d)", req.Header["Accept"], len(req.Header["Accept"]))
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
	user := UserFromStaff(staff)
	user.Refresh()
	ctx.Session.Values[kUserOL] = user
	ctx.Session.Values[kLastUid] = staff.Uid

	res["ok"] = true
	res["referer"] = ctx.Referer
	return outputJson(res, w)
	// http.Redirect(w, req, reverse("welcome"), http.StatusSeeOther)
}

func logout(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	delete(ctx.Session.Values, kUserOL)
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

func profileForm(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	if ctx.User == nil {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}
	staff, err := backends.GetStaff(ctx.User.Uid)
	if err != nil {
		return err
	}

	return T("profile.html").Execute(w, map[string]interface{}{
		"ctx":   ctx,
		"staff": staff,
	})
}

func profilePost(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	if ctx.User == nil {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}
	res := make(osin.ResponseData)
	// filed, value := req.PostFormValue("name"), req.PostFormValue("value")
	values := make(map[string]string)
	for input, field := range models.ProfileEditables {
		value := req.PostFormValue(input)
		if value != "" {
			values[field] = value
		}
	}
	password := req.PostFormValue("password")
	err := backends.ProfileModify(ctx.User.Uid, password, values)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
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
	w.Header().Set("Expires", "Fri, 02 Oct 1998 20:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate")

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

	ctx.afterHandle()

	//save the session
	if len(ctx.Session.Values) > 0 { // session not empty only
		if err = ctx.Session.Save(req, buf); err != nil {
			log.Printf("session.save error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//apply the buffered response to the writer
	buf.Apply(w)
}

func debugf(format string, args ...interface{}) {
	if Settings.Debug {
		log.Printf(format, args...)
	}
}
