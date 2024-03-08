package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-osin/osin"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/oauth"
	"github.com/liut/staffio/pkg/models/oidc"
	"github.com/liut/staffio/pkg/settings"
	"github.com/liut/staffio/pkg/web/i18n"
)

func (s *server) oauth2AuthorizeFirst(c *gin.Context, ar *osin.AuthorizeRequest) {
	s.Render(c, "authorize.html", map[string]interface{}{
		"link":          c.Request.RequestURI,
		"response_type": ar.Type,
		"scopes":        s.loadScopes(c.Request),
		"client":        ar.Client.(*oauth.Client),
		"ctx":           c,
	})
}

// Authorization code endpoint
func (s *server) oauth2Authorize(c *gin.Context) {
	resp := s.osvr.NewResponse()
	defer resp.Close()

	r := c.Request
	user := UserWithContext(c)
	store := s.service.OSIN()

	if ar := s.osvr.HandleAuthorizeRequest(resp, r); ar != nil {
		logger().Debugw("HandleAuthorizeRequest", "client", ar.Client)
		isAuthorized := store.IsAuthorized(ar.Client.GetId(), user.UID)
		if !isAuthorized && r.Method == "GET" {
			s.oauth2AuthorizeFirst(c, ar)
			return
		}

		if c.PostForm("authorize") == "1" {
			isAuthorized = true
		}

		if isAuthorized {

			// These values would be tied to the end user authorizing the client.
			err := s.oauth2UserData(c, ar, user)
			if err != nil {
				resp.SetError("get_user_error", "staff not found")
				resp.InternalError = err
			} else {
				ar.Authorized = true
			}

			s.osvr.FinishAuthorizeRequest(resp, r, ar)
			if r.PostForm.Get("remember") != "" {
				err := store.SaveAuthorized(ar.Client.GetId(), user.UID)
				if err != nil {
					logger().Infow("SaveAuthorized fail", "err", err)
				}
			}
		} else {
			resp.SetRedirect("/")
		}

	}

	if resp.IsError && resp.InternalError != nil {
		logger().Infow("authorize fail", "eid", resp.ErrorId, "err", resp.InternalError)
	}

	logger().Debugw("oauthAuthorize", "resp", resp)
	osin.OutputJSON(resp, c.Writer, r)
}

func (s *server) oauth2UserData(c *gin.Context, ar *osin.AuthorizeRequest,
	user *User) error {

	staff, err := s.service.Get(user.UID)
	if err != nil {
		return err
	}
	scopes := make(map[string]bool)
	for _, s := range strings.Fields(ar.Scope) {
		scopes[s] = true
	}
	// If the "openid" connect scope is specified, attach an ID Token to the
	// authorization response.
	//
	// The ID Token will be serialized and signed during the code for token exchange.
	if scopes["openid"] {
		now := time.Now()
		idToken := IDToken{
			Issuer:     settings.Current.BaseURL,
			UserID:     staff.UID,
			ClientID:   ar.Client.GetId(),
			Expiration: now.Add(time.Hour).Unix(),
			IssuedAt:   now.Unix(),
			Nonce:      c.Query("nonce"),
			UID:        staff.UID,
		}

		if scopes["profile"] {
			idToken.Name = staff.Name()
			idToken.GivenName = staff.GivenName
			idToken.FamilyName = staff.Surname
			idToken.BirthDate = staff.Birthday
			idToken.Nickname = staff.Nickname
			idToken.Locale = staff.OrgDepartment
		}

		if scopes["email"] {
			t := true
			idToken.Email = staff.Email
			idToken.EmailVerified = &t
		}
		// NOTE: The storage must be able to encode and decode this object.
		ar.UserData = &idToken
	} else {
		ar.UserData = oauth.JSONKV{"uid": user.UID}
	}

	return nil
}

func (s *server) buildJWT(staff *models.Staff, client, scope string) (string, error) {
	if strings.Contains(scope, "openid") {
		now := time.Now()
		idToken := &IDToken{
			Issuer:     settings.Current.BaseURL,
			UserID:     staff.UID,
			ClientID:   client,
			Expiration: now.Add(time.Hour).Unix(),
			IssuedAt:   now.Unix(),
			UID:        staff.UID,
		}
		return s.tkgen.GenerateIDToken(idToken)
	}
	return "", fmt.Errorf("invalid scope: %s", scope)
}

func (s *server) loadScopes(r *http.Request) (data []oauth.Scope) {
	p := i18n.GetPrinter(r)
	for _, s := range strings.Fields(r.FormValue("scope")) {
		scope := i18n.Scope(s)
		if !scope.Valid() {
			continue
		}
		data = append(data, oauth.Scope{
			Label:       scope.LabelP(p),
			Description: scope.DescriptionP(p),
		})
	}
	return
}

// Access token endpoint
func (s *server) oauth2Token(c *gin.Context) {
	resp := s.osvr.NewResponse()
	defer resp.Close()
	r := c.Request

	var (
		uid   string
		user  *User
		staff *models.Staff
		err   error
	)
	if ar := s.osvr.HandleAccessRequest(resp, r); ar != nil {
		logger().Debugw("HandleAccessRequest", "code", ar.Code, "scope", ar.Scope)
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			kv, _ := oauth.ToJSONKV(ar.UserData)
			if v, ok := kv["uid"]; ok {
				uid = v.(string)
			}

			staff, err = s.service.Get(uid)
			if err != nil {
				resp.SetError("get_user_error", "staff not found")
				resp.InternalError = err
			} else {
				user = UserFromStaff(staff)
			}
			if idt, err := s.buildJWT(staff, ar.Client.GetId(), ar.Scope); err == nil {
				resp.Output["id_token"] = idt
			}
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.UserData = nil
			// TODO: load refresh
			ar.Authorized = true
		case osin.PASSWORD:
			var staff *models.Staff
			if staff, err = s.service.Authenticate(ar.Username, ar.Password); err != nil {
				resp.SetError("authentication_failed", err.Error())
				break
			}
			ar.Authorized = true
			ar.UserData = oauth.JSONKV{"uid": staff.UID}
			user = UserFromStaff(staff)

		case osin.CLIENT_CREDENTIALS:
			ar.UserData = nil
			ar.Authorized = true
		case osin.ASSERTION:
			ar.UserData = nil
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}
		s.osvr.FinishAccessRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		logger().Infow("token ERROR", "err", resp.InternalError)
	}
	if !resp.IsError {
		if uid != "" {
			resp.Output["uid"] = uid
			resp.Output["is_keeper"] = s.IsKeeper(uid)
		}
		if user != nil {
			resp.Output["user"] = user
		}

	}

	logger().Infow("oauthToken", "resp", resp)

	osin.OutputJSON(resp, c.Writer, r)
}

// Information endpoint
func (s *server) oauth2Info(c *gin.Context) {
	resp := s.osvr.NewResponse()
	defer resp.Close()
	r := c.Request

	if ir := s.osvr.HandleInfoRequest(resp, r); ir != nil {
		logger().Debugw("HandleInfoRequest", "code", ir.Code, "accessToken", ir.AccessData.AccessToken)
		var (
			uid   string
			topic = c.Param("topic")
		)
		logger().Infow("param", "topic", topic)
		kv, _ := oauth.ToJSONKV(ir.AccessData.UserData)
		if v, ok := kv["uid"]; ok {
			uid = v.(string)
		}
		staff, err := s.service.Get(uid)
		if err != nil {
			resp.SetError("get_user_error", "staff not found")
			resp.InternalError = err
		} else {
			resp.Output["uid"] = uid
			if strings.HasPrefix(topic, "me") {
				resp.Output["me"] = staff
				if len(topic) > 3 && topic[2] == '+' {
					if arr := strings.Split(topic[3:], "+"); len(arr) > 0 {
						logger().Infow("search groups", "arr", arr)
						gm := make(map[string]interface{})
						for _, gn := range arr {
							if gn == "" {
								continue
							}
							gm[gn] = s.InGroup(gn, uid)
						}
						logger().Infow("result", "gm", gm)
						resp.Output["group"] = gm
					}
				} else if len(topic) > 3 && topic[2] == '|' {
					if arr := strings.Split(topic[3:], "|"); len(arr) > 0 {
						var roles []string
						for _, gn := range arr {
							if s.InGroup(gn, uid) {
								roles = append(roles, gn)
							}
						}
						logger().Infow("result", "roles", roles)
						resp.Output["group"] = roles
					}
				}

			} else if topic == "staff" {
				resp.Output["staff"] = staff
			} else if topic == "grafana" || topic == "generic" {
				resp.Output["name"] = staff.GetName()
				resp.Output["login"] = staff.UID
				resp.Output["sub"] = staff.UID
				resp.Output["preferred_username"] = staff.UID
				resp.Output["username"] = staff.UID
				resp.Output["email"] = staff.Email
				resp.Output["attributes"] = map[string][]string{} // TODO: fill attributes
			}

		}
		s.osvr.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError && resp.InternalError != nil {
		logger().Infow("info ERROR", "err", resp.InternalError)
	}

	osin.OutputJSON(resp, c.Writer, r)
}

func (s *server) oidcDiscovery(c *gin.Context) {
	r := c.Request
	baseURL := fmt.Sprintf("%s://%s", RequestScheme(r), r.Host)
	od := oidc.DiscoveryWith(baseURL)
	c.JSON(200, &od)
}

func RequestScheme(r *http.Request) string {
	if s := r.URL.Scheme; s != "" {
		return s
	}
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		return s
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
