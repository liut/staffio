package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RangelReale/osin"

	"lcgc/platform/staffio/backends"
	"lcgc/platform/staffio/models"
	. "lcgc/platform/staffio/settings"
)

func welcome(ctx *Context) (err error) {

	if Settings.Debug {
		log.Printf("session Name: %s, Values: %d", ctx.Session.Name(), len(ctx.Session.Values))
		log.Printf("ctx User %v", ctx.User)
	}

	//execute the template
	return ctx.Render("welcome.html", map[string]interface{}{
		"ctx": ctx,
	})
}

func loginForm(ctx *Context) (err error) {
	return ctx.Render("login.html", map[string]interface{}{
		"ctx": ctx,
	})
}

func login(ctx *Context) error {
	req := ctx.Request
	uid, password := req.FormValue("username"), req.FormValue("password")
	// log.Printf("accept: %v (%d)", req.Header["Accept"], len(req.Header["Accept"]))
	res := make(osin.ResponseData)
	if !backends.Authenticate(uid, password) {
		// ctx.Session.AddFlash("Invalid Username/Password")

		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "password"}
		return outputJson(res, ctx.Writer)
	}

	staff, err := backends.GetStaff(uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Load user failed"}
		return outputJson(res, ctx.Writer)
	}

	//store the user id in the values and redirect to welcome
	user := UserFromStaff(staff)
	user.Refresh()
	ctx.Session.Values[kUserOL] = user
	ctx.Session.Values[kLastUid] = staff.Uid

	res["ok"] = true
	res["referer"] = ctx.Referer
	return outputJson(res, ctx.Writer)
	// http.Redirect(ctx.Writer, req, reverse("welcome"), http.StatusSeeOther)
}

func logout(ctx *Context) error {
	delete(ctx.Session.Values, kUserOL)
	http.Redirect(ctx.Writer, ctx.Request, reverse("welcome"), http.StatusSeeOther)
	return nil
}

func passwordForm(ctx *Context) error {
	return ctx.Render("password.html", map[string]interface{}{
		"ctx": ctx,
	})
}

func passwordChange(ctx *Context) error {
	req := ctx.Request
	uid, pwdOld, pwdNew := req.FormValue("username"), req.FormValue("old_password"), req.FormValue("new_password")
	res := make(osin.ResponseData)
	if !backends.Authenticate(uid, pwdOld) {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "old_password"}
		return outputJson(res, ctx.Writer)
	}
	err := backends.PasswordChange(uid, pwdOld, pwdNew)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "old_password"}
	} else {
		res["ok"] = true
	}

	return outputJson(res, ctx.Writer)
}

func profileForm(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	staff, err := backends.GetStaff(ctx.User.Uid)
	if err != nil {
		return err
	}

	return ctx.Render("profile.html", map[string]interface{}{
		"ctx":   ctx,
		"staff": staff,
	})
}

func profilePost(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	res := make(osin.ResponseData)
	req := ctx.Request
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

	return outputJson(res, ctx.Writer)
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

func debugf(format string, args ...interface{}) {
	if Settings.Debug {
		log.Printf(format, args...)
	}
}
