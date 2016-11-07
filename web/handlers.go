package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin/binding"

	"lcgc/platform/staffio/backends"
	"lcgc/platform/staffio/backends/exmail"
	"lcgc/platform/staffio/models"
	. "lcgc/platform/staffio/settings"
)

func clientsForm(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
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
	return ctx.Render("clients.html", map[string]interface{}{
		"ctx":     ctx,
		"clients": clients,
	})
}

func clientsPost(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}
	res := make(osin.ResponseData)
	req := ctx.Request
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
			return outputJson(res, ctx.Writer)
		}

	} else {

		pk, name, value := req.PostFormValue("pk"), req.PostFormValue("name"), req.PostFormValue("value")
		log.Printf("clientsPost: pk %s, name %s, value %s", pk, name, value)
		if pk == "" {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is empty"}
			return outputJson(res, ctx.Writer)
		}
		// id, err := strconv.ParseUint(pk, 10, 32)
		client, err = backends.GetClientWithCode(pk)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is invalid or not found"}
			return outputJson(res, ctx.Writer)
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
			return outputJson(res, ctx.Writer)
		}
		res["ok"] = true
		res["id"] = client.Id
		return outputJson(res, ctx.Writer)
	}

	res["ok"] = false
	res["error"] = map[string]string{"message": "invalid operation"}
	return outputJson(res, ctx.Writer)
}

func scopesForm(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}
	scopes, err := backends.LoadScopes()
	if err != nil {
		return err
	}
	return ctx.Render("scopes.html", map[string]interface{}{
		"ctx":    ctx,
		"scopes": scopes,
	})
}

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

func contactsTable(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}

	staffs := backends.LoadStaffs()
	models.ByUid.Sort(staffs)

	return ctx.Render("contact.html", map[string]interface{}{
		"staffs": staffs,
		"ctx":    ctx,
	})
}

func staffForm(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}

	var (
		inEdit bool
		uid    = ctx.Vars["uid"]
		staff  *models.Staff
		data   = map[string]interface{}{
			"ctx": ctx,
		}
	)

	if uid != "" && uid != "new" {
		inEdit = true
		staff, err = backends.GetStaff(uid)
		if err != nil {
			return
		}
		data["staff"] = staff
	}
	data["inEdit"] = inEdit
	return ctx.Render("staff_edit.html", data)
}

func staffPost(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return
	}
	req := ctx.Request

	var (
		uid           = ctx.Vars["uid"]
		estaff, staff *models.Staff
		res           = make(osin.ResponseData)
		op            = req.FormValue("op")
	)
	if uid == "" || uid == "new" {
		uid = req.PostFormValue("uid")
	}

	if uid == "" || uid == "new" {
		return fmt.Errorf("empty uid")
	} else {
		estaff, err = backends.GetStaff(uid)
		if err != nil {
			log.Printf("backends.GetStaff err %s", err)
			estaff = nil
		}
	}

	email := uid + "@" + Settings.EmailDomain
	if op == "fetch-exmail" && uid != "" {
		staff, err = exmail.GetStaff(email)
		if err != nil {
			log.Printf("GetStaff err %s", err)
			return err
		}
		log.Print(staff)
		if estaff != nil {
			staff.CommonName = estaff.CommonName
			staff.Surname = estaff.Surname
			staff.GivenName = estaff.GivenName
			staff.Gender = estaff.Gender
			staff.Nickname = estaff.Nickname
			if estaff.Mobile != "" {
				staff.Mobile = estaff.Mobile
			}
			if estaff.Email != "" {
				staff.Email = estaff.Email
			}
			if estaff.EmployeeNumber != "" {
				staff.EmployeeNumber = estaff.EmployeeNumber
			}
			if estaff.EmployeeType != "" {
				staff.EmployeeType = estaff.EmployeeType
			}
			if estaff.Description != "" {
				staff.Description = estaff.Description
			}
		}
		res["ok"] = true
		res["staff"] = staff
		outputJson(res, ctx.Writer)
	} else if op == "store" {
		fb := binding.Form
		staff = new(models.Staff)
		err = fb.Bind(req, staff)
		if err != nil {
			log.Printf("bind %v: %s", staff, err)
			return
		}
		log.Print(staff)

		// sn, gn := req.PostFormValue("sn"), req.PostFormValue("sn")
		// cn := sn + gn
		// staff = models.NewStaff(uid, cn, email)
		// staff.Surname = sn
		// staff.GivenName = gn
		// staff.Mobile = req.PostFormValue("mobile")
		err = backends.StoreStaff(staff)
		if err == nil {
			res["ok"] = true
			res["referer"] = reverse("contacts")
			outputJson(res, ctx.Writer)
		}
	}

	return
}

func staffDelete(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return
	}

	var (
		uid = ctx.Vars["uid"]
		res = make(osin.ResponseData)
	)

	if uid == "" || uid == "new" {
		return fmt.Errorf("empty uid")
	}

	if uid == ctx.User.Uid {
		res["ok"] = false
		res["error"] = "Can not delete yourself"
		return outputJson(res, ctx.Writer)
	}

	_, err = backends.GetStaff(uid)
	if err != nil {
		log.Printf("backends.GetStaff err %s", err)
		http.Error(ctx.Writer, err.Error(), http.StatusNotFound)
		return err
	}
	err = backends.DeleteStaff(uid)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return err
	}

	res["ok"] = true
	return outputJson(res, ctx.Writer)

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
