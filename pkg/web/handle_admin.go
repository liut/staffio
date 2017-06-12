package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin/binding"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/models"
	. "lcgc/platform/staffio/pkg/settings"
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
		// log.Printf("clientsPost: pk %s, name %s, value %s", pk, name, value)
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
		// log.Print(staff, uint8(staff.Gender))
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
		staff, err = backends.GetStaffFromExmail(email)
		if err != nil {
			log.Printf("GetStaff err %s", err)
			return err
		}
		// log.Print(staff)
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
