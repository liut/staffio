package web

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/RangelReale/osin"

	"lcgc/platform/staffio/backends"
	"lcgc/platform/staffio/models"
	. "lcgc/platform/staffio/settings"
)

// Authorization code endpoint
func oauthAuthorize(ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()

	r := ctx.Request

	if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
		link := fmt.Sprintf("/authorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
			ar.Type, ar.Client.GetId(), url.QueryEscape(ar.RedirectUri), ar.State, ar.Scope)
		// HANDLE LOGIN PAGE HERE
		if ctx.User == nil {
			ctx.Referer = link
			return loginForm(ctx)
			// resp.SetRedirect(reverse("login") + "?referer=" + reverse("authorize"))
		} else {
			if r.Method == "GET" {
				scopes, err := backends.LoadScopes()
				if err != nil {
					return err
				}
				return ctx.Render("authorize.html", map[string]interface{}{
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
	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}

// Access token endpoint
func oauthToken(ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()
	r := ctx.Request

	var (
		uid   string = ""
		user  *User
		staff *models.Staff
	)
	if ar := server.HandleAccessRequest(resp, r); ar != nil {
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

	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}

// Information endpoint
func oauthInfo(ctx *Context) (err error) {
	resp := server.NewResponse()
	defer resp.Close()
	r := ctx.Request

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
					gn := topic[3:]
					resp.Output[gn] = InGroup(gn, uid)
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

	osin.OutputJSON(resp, ctx.Writer, r)
	return resp.InternalError
}
