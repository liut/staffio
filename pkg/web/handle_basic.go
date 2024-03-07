package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-osin/osin"

	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/cas"
	"github.com/liut/staffio/pkg/web/apis"
	"github.com/liut/staffio/pkg/web/i18n"
)

func (s *server) loginForm(c *gin.Context) {
	service := c.Request.FormValue("service")
	tgc := GetTGC(c)
	if service != "" && tgc != nil {
		st := cas.NewTicket("ST", service, tgc.UID, false)
		err := s.service.SaveTicket(st)
		if err != nil {
			return
		}
		c.Redirect(302, service+"?ticket="+st.Value)
		return
	}
	s.Render(c, "login.html", map[string]interface{}{
		"ctx":     c,
		"service": service,
	})
}

// loginPost ...
// @Tag staffio
// @Summary login
// @Description login
// @ID api-1-login-post
// @Accept  x-www-form-urlencoded,mpfd,json
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} apis.RespDone
// @Failure 400 {object} apis.RespFail
// @Failure 401 {object} apis.RespFail
// @Failure 500 {object} apis.RespFail
// @Router /api/login [post]
func (s *server) loginPost(c *gin.Context) {
	var param loginParam
	res := make(osin.ResponseData)
	if err := c.Bind(&param); err != nil {
		apis.Fail(c, 400, err)
		return
	}
	// req := c.Request
	// uid, password := req.PostFormValue("username"), req.PostFormValue("password")

	var (
		staff *models.Staff
		err   error
	)
	if staff, err = s.service.Authenticate(param.Username, param.Password); err != nil {
		apis.Fail(c, 401, i18n.ErrLoginFailed, "password")
		return
	}

	//store the user id in the values and redirect to welcome
	signinStaffGin(c, staff)
	referer := param.Referer
	res["ok"] = true
	if param.Service != "" {
		st := cas.NewTicket("ST", param.Service, staff.UID, true)
		err = s.service.SaveTicket(st)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		NewTGC(c, st)
		res["referer"] = param.Service + "?ticket=" + st.Value
		log.Printf("ref: %q", res["referer"])
	} else {
		if referer == "" {
			referer = "/"
		}
		res["referer"] = referer
	}
	res["status"] = 0
	c.JSON(http.StatusOK, res)
	// http.Redirect(c.Writer, req, reverse("welcome"), http.StatusSeeOther)
}

// for staff/verify
func (s *server) me(c *gin.Context) {
	user, err := authzr.UserFromRequest(c.Request)
	if err != nil {
		apiError(c, 1, nil)
		return
	}
	team, err := s.service.Team().GetWithMember(user.UID)
	if err == nil {
		user.TeamID = int64(team.ID)
	} else {
		log.Printf("get team with member %s ERR %s", user.UID, err)
	}
	if s.IsKeeper(user.UID) {
		user.Roles = append(user.Roles, "admin")
	}
	if team.Leaders.Contains(user.UID) {
		user.Roles = append(user.Roles, "leader")
	}
	user.Watchings = s.service.Watch().Gets(user.UID).UIDs()
	apiOk(c, user, 0)
}

func (s *server) logout(c *gin.Context) {
	authzr.Signout(c.Writer)
	DeleteTGC(c)
	if IsAjax(c.Request) {
		apiOk(c, true, 0)
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (s *server) passwordForm(c *gin.Context) {
	s.Render(c, "password.html", map[string]interface{}{
		"ctx": c,
	})
}

// passwordChange ...
// @Tag staffio
// @Summary Change password
// @Description change password
// @ID api-1-password-post
// @Accept  x-www-form-urlencoded,mpfd,json
// @Produce  json
// @Param old_password formData string true "Old Password"
// @Param new_password formData string true "New Password"
// @Param password_confirm formData string true "Confirm Password"
// @Success 200 {object} apis.RespDone
// @Failure 400 {object} apis.RespFail
// @Failure 401 {object} apis.RespFail
// @Failure 500 {object} apis.RespFail
// @Router /api/password [post]
func (s *server) passwordChange(c *gin.Context) {
	var param passwordParam
	res := make(osin.ResponseData)
	if err := c.Bind(&param); err != nil {
		res["ok"] = false
		res["error"] = err.Error()
		res["status"] = ERROR_PARAM
		c.JSON(400, res)
		return
	}
	user := UserWithContext(c)
	if _, err := s.service.Authenticate(user.UID, param.OldPassword); err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username/Password", "field": "old_password"}
		res["status"] = ERROR_PARAM
		c.JSON(http.StatusOK, res)
		return
	}
	err := s.service.PasswordChange(user.UID, param.OldPassword, param.NewPassword)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "old_password"}
		res["status"] = ERROR_DB
	} else {
		res["ok"] = true
		res["status"] = 0
	}

	c.JSON(http.StatusOK, res)
}

func (s *server) passwordForgotForm(c *gin.Context) {
	s.Render(c, "password_forgot.html", map[string]interface{}{
		"ctx": c,
	})
}

// passwordForgot ...
// @Tag staffio
// @Summary Forgot password
// @Description forgot password
// @ID api-1-password-forgot-post
// @Accept  x-www-form-urlencoded,mpfd,json
// @Produce  json
// @Param username formData string true "Login name"
// @Param mobile formData string true "Mobile number"
// @Param email formData string true "Email address"
// @Success 200 {object} apis.RespDone
// @Failure 400 {object} apis.RespFail
// @Failure 401 {object} apis.RespFail
// @Failure 500 {object} apis.RespFail
// @Router /api/password/forgot [post]
func (s *server) passwordForgot(c *gin.Context) {
	var param forgotParam
	res := make(osin.ResponseData)
	if err := c.Bind(&param); err != nil {
		res["ok"] = false
		res["error"] = err.Error()
		res["status"] = ERROR_PARAM
		c.JSON(400, res)
		return
	}

	staff, err := s.service.Get(param.Username)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username", "field": "username"}
		c.JSON(http.StatusOK, res)
		return
	}
	if staff.Email != param.Email {
		res["ok"] = false
		res["error"] = map[string]string{"message": "No such email address", "field": "email"}
		c.JSON(http.StatusOK, res)
		return
	}
	// if staff.Mobile != param.Mobile {
	// 	res["ok"] = false
	// 	res["error"] = map[string]string{"message": "The mobile number is a mismatch", "field": "mobile"}
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }
	logger().Infow("forgot", "req.host", c.Request.Host, "url.host", c.Request.URL.Host)
	err = s.service.PasswordForgot(ContextWithSiteFromRequest(c.Request), common.AtEmail, param.Email, param.Username)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "username"}
	} else {
		res["ok"] = true
		res["status"] = 0
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
	s.Render(c, "password_reset.html", map[string]interface{}{
		"ctx":   c,
		"token": token,
		"uid":   uid,
	})
}

// passwordReset ...
// @Tag staffio
// @Summary Reset password
// @Description reset password, form:rt, json:token
// @ID api-1-password-reset-post
// @Accept  x-www-form-urlencoded,mpfd,json
// @Produce  json
// @Param username formData string true "Login name"
// @Param password formData string true "Password"
// @Param password_confirm formData string true "Confirm Password"
// @Param rt formData string true "Token"
// @Success 200 {object} apis.RespDone
// @Failure 400 {object} apis.RespFail
// @Failure 401 {object} apis.RespFail
// @Failure 500 {object} apis.RespFail
// @Router /api/password/reset [post]
func (s *server) passwordReset(c *gin.Context) {
	var param resetParam
	res := make(osin.ResponseData)
	if err := c.Bind(&param); err != nil {
		res["ok"] = false
		res["error"] = err.Error()
		res["status"] = ERROR_PARAM
		c.JSON(400, res)
		return
	}

	if param.Password != param.Password2 {
		res["ok"] = false
		res["error"] = map[string]string{"message": "Invalid Username or Password", "field": "password"}
		res["status"] = ERROR_PARAM
		c.JSON(http.StatusOK, res)
		return
	}
	err := s.service.PasswordResetWithToken(param.Username, param.Token, param.Password)
	if err != nil {
		res["ok"] = false
		res["status"] = ERROR_DB
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
	} else {
		res["ok"] = true
		res["status"] = 0
	}
	c.JSON(http.StatusOK, res)
}

func (s *server) profileForm(c *gin.Context) {
	user := UserWithContext(c)
	staff, err := s.service.Get(user.UID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	s.Render(c, "profile.html", map[string]interface{}{
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
	err = s.service.ProfileModify(user.UID, password, staff)
	if err != nil {
		res["ok"] = false
		res["error"] = map[string]string{"message": err.Error(), "field": "password"}
	} else {
		res["ok"] = true
	}

	c.JSON(http.StatusOK, res)
}
