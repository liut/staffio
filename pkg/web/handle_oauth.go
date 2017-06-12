package web

import (
	// "fmt"
	"log"
	// "net/url"
	"strings"

	"github.com/RangelReale/osin"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/models"
	. "lcgc/platform/staffio/pkg/settings"
)

// Authorization code endpoint
func oauthAuthorize(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	resp := ws.osvr.NewResponse()
	defer resp.Close()

	r := ctx.Request

	if ar := ws.osvr.HandleAuthorizeRequest(resp, r); ar != nil {
		// link := fmt.Sprintf("/authorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		// 	ar.Type, ar.Client.GetId(), url.QueryEscape(ar.RedirectUri), ar.State, ar.Scope)
		if backends.IsAuthorized(ar.Client.GetId(), ctx.User.Uid) {
			ar.UserData = ctx.User.Uid
			ar.Authorized = true
			ws.osvr.FinishAuthorizeRequest(resp, r, ar)
		} else {
			if r.Method == "GET" {
				scopes, err := backends.LoadScopes()
				if err != nil {
					return err
				}
				return ctx.Render("authorize.html", map[string]interface{}{
					"link":          r.RequestURI,
					"response_type": ar.Type,
					"scopes":        scopes,
					"client":        ar.Client.(*models.Client),
					"ctx":           ctx,
				})
			}

			if r.PostForm.Get("authorize") == "1" {
				ar.UserData = ctx.User.Uid
				ar.Authorized = true
				ws.osvr.FinishAuthorizeRequest(resp, r, ar)
				if r.PostForm.Get("remember") != "" {
					err := backends.SaveAuthorized(ar.Client.GetId(), ctx.User.Uid)
					if err != nil {
						log.Printf("remember ERR %s", err)
					}
				}
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
	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}

// Access token endpoint
func oauthToken(ctx *Context) (err error) {
	resp := ws.osvr.NewResponse()
	defer resp.Close()
	r := ctx.Request

	var (
		uid   string = ""
		user  *User
		staff *models.Staff
	)
	if ar := ws.osvr.HandleAccessRequest(resp, r); ar != nil {
		debugf("ar Code %s Scope %s", ar.Code, ar.Scope)
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			uid = ar.UserData.(string)
			staff, err = backends.GetStaff(uid)
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
		ws.osvr.FinishAccessRequest(resp, r, ar)
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

	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}

// Information endpoint
func oauthInfo(ctx *Context) (err error) {
	resp := ws.osvr.NewResponse()
	defer resp.Close()
	r := ctx.Request

	if ir := ws.osvr.HandleInfoRequest(resp, r); ir != nil {
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
					gn := topic[3:]
					resp.Output[gn] = InGroup(gn, uid)
				}
			} else if topic == "staff" {
				resp.Output["staff"] = staff
			} else if topic == "grafana" || topic == "generic" {
				resp.Output["name"] = staff.Name()
				resp.Output["login"] = staff.Uid
				resp.Output["username"] = staff.Uid
				resp.Output["email"] = staff.Email
				resp.Output["attributes"] = map[string][]string{} // TODO: fill attributes
			}

		}
		ws.osvr.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError && resp.InternalError != nil {
		log.Printf("info ERROR: %s\n", resp.InternalError)
	}

	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}
