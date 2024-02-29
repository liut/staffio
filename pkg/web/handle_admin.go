package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-osin/osin"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/backends/qqexmail"
	"github.com/liut/staffio/pkg/backends/wechatwork"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/oauth"
	"github.com/liut/staffio/pkg/settings"
)

func (s *server) clientsGet(c *gin.Context) {
	spec := new(oauth.ClientSpec)
	if err := c.Bind(spec); err != nil {
		apiError(c, 400, err)
		return
	}

	clients, err := s.service.OSIN().LoadClients(spec)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if IsAjax(c.Request) {
		apiOk(c, clients, len(clients))
		return
	}
	s.Render(c, "clients.html", map[string]interface{}{
		"ctx":     c,
		"clients": clients,
	})
}

type inlineEdit struct {
	PK    string `form:"pk" json:"pk" binding:"required" description:"主键"`
	Field string `form:"name" json:"name" binding:"required" description:"名称"`
	Value string `form:"value" json:"value" binding:"required"`
}

type clientParam struct {
	ID          string `form:"id" json:"id" binding:"required"`
	Secret      string `form:"secret" json:"secret"`
	RedirectURI string `form:"redirect_uri" json:"redirect_uri"`
	Name        string `form:"name" json:"name"`
}

func (s *server) clientsPost(c *gin.Context) {
	res := make(osin.ResponseData)
	req := c.Request
	var (
		client *oauth.Client
		err    error
	)

	if req.FormValue("op") == "new" {
		// create new client
		client = backends.GenNewClient(req.PostFormValue("name"), req.PostFormValue("redirect_uri"))
		if err = s.service.OSIN().SaveClient(client); err != nil {
			apiError(c, 1, err)
			return
		}
		res["ok"] = true
		res["id"] = client.ID
		c.JSON(http.StatusOK, res)
		return
	}
	var (
		inline inlineEdit
		param  clientParam
	)
	if req.FormValue("pk") != "" {
		err = c.Bind(&inline)
		if err != nil {
			apiError(c, 400, err)
			return
		}
		client, err = s.service.OSIN().LoadClient(inline.PK)
		if err != nil {
			logger().Warnw("get client ERR ", "inline", inline, "err", err)
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is invalid or not found"}
			c.JSON(http.StatusOK, res)
			return
		}
		switch inline.Field {
		case "name":
			client.Meta.Name = inline.Value
		case "secret":
			client.Secret = inline.Value
		case "redirect_uri":
			client.RedirectURI = inline.Value
		default:
			logger().Infow("invalid", "field", inline.Field)
			apiError(c, 400, "invalid field")
			return
		}
	} else if err = c.Bind(&param); err == nil {
		client, err = s.service.OSIN().LoadClient(param.ID)
		if len(param.Name) > 0 && client.Meta.Name != param.Name {
			client.Meta.Name = param.Name
		}
		if len(param.Secret) > 0 && client.Secret != param.Secret {
			client.Secret = param.Secret
		}
		if len(param.RedirectURI) > 0 && client.RedirectURI != param.RedirectURI {
			client.RedirectURI = param.RedirectURI
		}
	} else {
		logger().Warnw("bind failed ", c.Request.Method, c.Request.RequestURI, "err", err)
		apiError(c, 400, err)
		return
	}

	if err == nil && client != nil {
		err = s.service.OSIN().SaveClient(client)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
			c.JSON(http.StatusOK, res)
			return
		}
		res["ok"] = true
		res["id"] = client.ID
		c.JSON(http.StatusOK, res)
		return
	}

	res["ok"] = false
	res["error"] = map[string]string{"message": "invalid operation"}
	c.JSON(http.StatusOK, res)
}

func (s *server) scopesForm(c *gin.Context) {
	scopes, err := s.service.OSIN().LoadScopes()
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	s.Render(c, "scopes.html", map[string]interface{}{
		"ctx":    c,
		"scopes": scopes,
	})
}

func (s *server) contactsTable(c *gin.Context) {
	var spec *models.Spec
	staffs := s.service.All(spec)
	models.ByUid.Sort(staffs)

	s.Render(c, "contact.html", map[string]interface{}{
		"staffs": staffs,
		"ctx":    c,
	})
}

func (s *server) staffForm(c *gin.Context) {

	var (
		inEdit bool
		uid    = c.Param("uid")
		staff  *models.Staff
		data   = map[string]interface{}{
			"ctx": c,
		}
		err error
	)

	if uid != "" && uid != "new" {
		inEdit = true
		staff, err = s.service.Get(uid)
		if err != nil {
			logger().Infow("get fail", "uid", uid, "err", err)
			return
		}
		// log.Print(staff, uint8(staff.Gender))
		data["staff"] = staff
	}
	data["inEdit"] = inEdit
	data["exmail"] = settings.Current.EmailCheck
	data["wxwork"] = settings.Current.WechatCorpID != ""
	s.Render(c, "staff_edit.html", data)
}

func (s *server) fetchEx(c *gin.Context) {
	src := c.Param("src")
	uid := c.Param("uid")

	estaff, err := s.service.Get(uid)
	if err != nil {
		logger().Infow("get fail", "uid", uid, "err", err)
		estaff = nil
	}
	logger().Debugw("fetch ex", "src", src, "uid", uid, "estaff", estaff)
	var staff *models.Staff
	switch src {
	case "wechat":
		exuser, err := s.wxAuth.GetUser(uid)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		staff = wechatwork.UserToStaff(exuser)
		logger().Infow("fetch-ex", "src", src, "user", exuser, "staff", staff)
	case "exmail":
		email := uid + "@" + settings.Current.EmailDomain
		staff, err = qqexmail.GetStaffFromExmail(email)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			logger().Infow("get fail", "uid", uid, "err", err)
			return
		}
	default:
		c.AbortWithError(400, err)
		return
	}

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
		if estaff.EmployeeNumber > 0 {
			staff.EmployeeNumber = estaff.EmployeeNumber
		}
		if estaff.EmployeeType != "" {
			staff.EmployeeType = estaff.EmployeeType
		}
		if estaff.Description != "" {
			staff.Description = estaff.Description
		}
	}
	res := make(osin.ResponseData)
	res["ok"] = true
	res["staff"] = staff
	c.JSON(http.StatusOK, res)
	return
}

func (s *server) staffPost(c *gin.Context) {
	req := c.Request
	uid := c.Param("uid")
	if uid == "" || uid == "new" {
		uid = req.PostFormValue("uid")
		if uid == "" || uid == "new" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	if req.FormValue("op") == "store" {
		fb := binding.Form
		staff := new(models.Staff)
		err := fb.Bind(req, staff)
		if err != nil {
			logger().Infow("bind fail", "staff", staff, "err", err)
			return
		}

		err = s.service.SaveStaff(staff)
		if err == nil {
			res := make(osin.ResponseData)
			res["ok"] = true
			res["referer"] = "/contacts"
			c.JSON(http.StatusOK, res)
		}
	}

	return
}

func (s *server) staffDelete(c *gin.Context) {

	var (
		uid = c.Param("uid")
		res = make(osin.ResponseData)
	)

	if uid == "" || uid == "new" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user := UserWithContext(c)

	if uid == user.UID {
		res["ok"] = false
		res["error"] = "Can not delete yourself"
		c.JSON(http.StatusOK, res)
		return
	}

	_, err := s.service.Get(uid)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		logger().Infow("get fail", "uid", uid, "err", err)
		return
	}
	err = s.service.Delete(uid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		logger().Infow("delete fail", "uid", uid, "err", err)
		return
	}

	res["ok"] = true
	c.JSON(http.StatusOK, res)

}

func (s *server) groupList(c *gin.Context) {

	data, _ := s.service.AllGroup()

	if strings.HasPrefix(c.Request.RequestURI, "/api/") || IsAjax(c.Request) {
		apiOk(c, data, len(data))
		return
	}

	s.Render(c, "group.html", map[string]interface{}{
		"groups": data,
		"ctx":    c,
	})
}

func (s *server) groupStore(c *gin.Context) {
	var group = new(backends.Group)
	if err := c.Bind(group); err != nil {
		apiError(c, http.StatusBadRequest, err)
		return
	}

	logger().Infow("groupStore", "group", group)

	if err := s.service.SaveGroup(group); err != nil {
		apiError(c, http.StatusServiceUnavailable, err)
		return
	}

	apiOk(c, true, 0)
}
