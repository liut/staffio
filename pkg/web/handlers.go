package web

import (
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	// "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"lcgc/platform/staffio/pkg/models"
	"lcgc/platform/staffio/pkg/models/cas"
	"lcgc/platform/staffio/pkg/models/common"
	. "lcgc/platform/staffio/pkg/settings"
)

func (s *server) loginForm(c *gin.Context) {
	service := c.Request.FormValue("service")
	tgc := GetTGC(c)
	if service != "" && tgc != nil {
		st := cas.NewTicket("ST", service, tgc.Uid, false)
		err := s.service.SaveTicket(st)
		if err != nil {
			return
		}
		c.Redirect(302, service+"?ticket="+st.Value)
		return
	}
	Render(c, "login.html", map[string]interface{}{
		"ctx":     c,
		"service": service,
	})
}

func (s *server) login(c *gin.Context) {
	req := c.Request
	// session := sessions.Default(c)
	uid, password := req.PostFormValue("username"), req.PostFormValue("password")
	service := req.FormValue("service")
	referer := req.PostFormValue("referer")
	res := make(osin.ResponseData)
	if err := s.service.Authenticate(uid, password); err != nil {

		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "password"}
		c.JSON(http.StatusOK, res)
		return
	}

	staff, err := s.service.Get(uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Load user failed"}
		c.JSON(http.StatusOK, res)
		return
	}

	//store the user id in the values and redirect to welcome
	user := UserFromStaff(staff)
	user.Refresh()
	user.toResponse(c.Writer)
	// session.Set(kUserOL, user)
	// session.Values[kUserOL] = user
	// session.Values[kLastUid] = staff.Uid
	res["ok"] = true
	if service != "" {
		st := cas.NewTicket("ST", service, user.Uid, true)
		err = s.service.SaveTicket(st)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		NewTGC(c, st)
		res["referer"] = service + "?ticket=" + st.Value
		log.Printf("ref: %q", res["referer"])
	} else {
		res["referer"] = referer
	}
	c.JSON(http.StatusOK, res)
	// http.Redirect(c.Writer, req, reverse("welcome"), http.StatusSeeOther)
}

func (s *server) logout(c *gin.Context) {
	DeleteTGC(c)
	c.Redirect(http.StatusSeeOther, "/")
}

func (s *server) passwordForm(c *gin.Context) {
	Render(c, "password.html", map[string]interface{}{
		"ctx": c,
	})
}

func (s *server) passwordChange(c *gin.Context) {
	req := c.Request
	uid, pwdOld, pwdNew := req.FormValue("username"), req.FormValue("old_password"), req.FormValue("new_password")
	res := make(osin.ResponseData)
	if err := s.service.Authenticate(uid, pwdOld); err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "old_password"}
		c.JSON(http.StatusOK, res)
		return
	}
	err := s.service.PasswordChange(uid, pwdOld, pwdNew)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "old_password"}
	} else {
		res["ok"] = true
	}

	c.JSON(http.StatusOK, res)
}

func (s *server) passwordForgotForm(c *gin.Context) {
	Render(c, "password_forgot.html", map[string]interface{}{
		"ctx": c,
	})
}

func (s *server) passwordForgot(c *gin.Context) {
	req := c.Request
	uid, email, mobile := req.FormValue("username"), req.FormValue("email"), req.FormValue("mobile")
	res := make(osin.ResponseData)
	staff, err := s.service.Get(uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username", "field": "username"}
		c.JSON(http.StatusOK, res)
		return
	}
	if staff.Email != email {
		res["ok"] = false
		res["error"] = map[string]string{"message": "No such email address", "field": "email"}
		c.JSON(http.StatusOK, res)
		return
	}
	if staff.Mobile != mobile {
		res["ok"] = false
		res["error"] = map[string]string{"message": "The mobile number is a mismatch", "field": "mobile"}
		c.JSON(http.StatusOK, res)
		return
	}
	err = s.service.PasswordForgot(common.AtEmail, email, uid)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "username"}
	} else {
		res["ok"] = true
	}
	c.JSON(http.StatusOK, res)
}

func (s *server) passwordResetForm(c *gin.Context) {
	req := c.Request

	token := req.FormValue("rt")
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uid, err := s.service.PasswordResetTokenVerify(token)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		// c.Halt(http.StatusBadRequest, fmt.Sprintf("Invalid Token: %s", err))
		return
	}
	Render(c, "password_reset.html", map[string]interface{}{
		"ctx":   c,
		"token": token,
		"uid":   uid,
	})
}

func (s *server) passwordReset(c *gin.Context) {
	req := c.Request

	token := req.FormValue("rt")
	if token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uid, passwd, passwd2 := req.FormValue("username"), req.FormValue("password"), req.FormValue("password_confirm")
	res := make(osin.ResponseData)
	if uid == "" || passwd != passwd2 {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username or Password", "field": "password"}
		c.JSON(http.StatusOK, res)
		return
	}
	err := s.service.PasswordResetWithToken(uid, token, passwd)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
	} else {
		res["ok"] = true
	}
	c.JSON(http.StatusOK, res)
}

func (s *server) profileForm(c *gin.Context) {
	user := UserWithContext(c)
	staff, err := s.service.Get(user.Uid)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	Render(c, "profile.html", map[string]interface{}{
		"ctx":   c,
		"staff": staff,
	})
}

func (s *server) profilePost(c *gin.Context) {
	user := UserWithContext(c)
	res := make(osin.ResponseData)
	req := c.Request
	password := req.PostFormValue("password")

	staff := new(models.Staff)
	err := binding.Form.Bind(req, staff)
	if err != nil {
		log.Printf("bind %v: %s", staff, err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = s.service.ProfileModify(user.Uid, password, staff)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
	} else {
		res["ok"] = true
	}

	c.JSON(http.StatusOK, res)
}

// func outputJson(res map[string]interface{}, w http.ResponseWriter) error {
// 	if w.Header().Get("Content-Type") == "" {
// 		w.Header().Set("Content-Type", "application/json")
// 	}

// 	encoder := json.NewEncoder(w)
// 	err := encoder.Encode(res)
// 	if err != nil {
// 		log.Printf("json encoding error: %s", err)
// 	}
// 	return err
// }

func debugf(format string, args ...interface{}) {
	if Settings.Debug {
		log.Printf(format, args...)
	}
}
