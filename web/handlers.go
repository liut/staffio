package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin/binding"

	"lcgc/platform/staffio/backends"
	"lcgc/platform/staffio/models"
	"lcgc/platform/staffio/models/cas"
	"lcgc/platform/staffio/models/common"
	. "lcgc/platform/staffio/settings"
)

func loginForm(ctx *Context) (err error) {
	service := ctx.Request.FormValue("service")
	tgc := GetTGC(ctx.Session)
	if service != "" && tgc != nil {
		st := cas.NewTicket("ST", service, tgc.Uid, false)
		err = backends.SaveTicket(st)
		if err != nil {
			return
		}
		ctx.Redirect(service + "?ticket=" + st.Value)
		return nil
	}
	return ctx.Render("login.html", map[string]interface{}{
		"ctx":     ctx,
		"service": service,
	})
}

func login(ctx *Context) error {
	req := ctx.Request
	uid, password := req.PostFormValue("username"), req.PostFormValue("password")
	service := req.FormValue("service")
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
	if service != "" {
		st := cas.NewTicket("ST", service, user.Uid, true)
		err = backends.SaveTicket(st)
		if err != nil {
			return err
		}
		NewTGC(ctx.Session, st)
		res["referer"] = service + "?ticket=" + st.Value
		log.Printf("ref: %q", res["referer"])
	} else {
		res["referer"] = ctx.Referer
	}
	return outputJson(res, ctx.Writer)
	// http.Redirect(ctx.Writer, req, reverse("welcome"), http.StatusSeeOther)
}

func logout(ctx *Context) error {
	delete(ctx.Session.Values, kUserOL)
	DeleteTGC(ctx.Session)
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

func passwordForgotForm(ctx *Context) error {
	return ctx.Render("password_forgot.html", map[string]interface{}{
		"ctx": ctx,
	})
}

func passwordForgot(ctx *Context) error {
	req := ctx.Request
	uid, email, mobile := req.FormValue("username"), req.FormValue("email"), req.FormValue("mobile")
	res := make(osin.ResponseData)
	staff, err := backends.GetStaff(uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username", "field": "username"}
		return outputJson(res, ctx.Writer)
	}
	if staff.Email != email {
		res["ok"] = false
		res["error"] = map[string]string{"message": "No such email address", "field": "email"}
		return outputJson(res, ctx.Writer)
	}
	if staff.Mobile != mobile {
		res["ok"] = false
		res["error"] = map[string]string{"message": "The mobile number is a mismatch", "field": "mobile"}
		return outputJson(res, ctx.Writer)
	}
	err = backends.PasswordForgot(common.AtEmail, email, uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "username"}
	} else {
		res["ok"] = true
	}
	return outputJson(res, ctx.Writer)
}

func passwordResetForm(ctx *Context) error {
	req := ctx.Request

	token := req.FormValue("rt")
	if token == "" {
		ctx.Halt(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return nil
	}

	uid, err := backends.PasswordResetTokenVerify(token)
	if err != nil {
		ctx.Halt(http.StatusBadRequest, fmt.Sprintf("Invalid Token: %s", err))
		return nil
	}
	return ctx.Render("password_reset.html", map[string]interface{}{
		"ctx":   ctx,
		"token": token,
		"uid":   uid,
	})
}

func passwordReset(ctx *Context) error {
	req := ctx.Request

	token := req.FormValue("rt")
	if token == "" {
		ctx.Halt(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return nil
	}
	uid, passwd, passwd2 := req.FormValue("username"), req.FormValue("password"), req.FormValue("password_confirm")
	res := make(osin.ResponseData)
	if uid == "" || passwd != passwd2 {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username or Password", "field": "password"}
		return outputJson(res, ctx.Writer)
	}
	err := backends.PasswordResetWithToken(uid, token, passwd)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
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
	password := req.PostFormValue("password")

	staff := new(models.Staff)
	err := binding.Form.Bind(req, staff)
	if err != nil {
		log.Printf("bind %v: %s", staff, err)
		return err
	}
	err = backends.ProfileModify(ctx.User.Uid, password, staff)
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
