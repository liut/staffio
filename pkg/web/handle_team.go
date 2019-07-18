package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/models/team"
)

type teamAddParam struct {
	Name string `json:"name" form:"name" binding:"required" valid:"[1:128]"`
}

type teamDeleteParam struct {
	ID int `json:"id" form:"form" binding:"required" valid:"required"`
}

func (s *server) teamListByRole(c *gin.Context) {
	role, err := strconv.Atoi(c.Request.FormValue("role"))
	if err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	data, err := s.service.Team().All(team.RoleType(role))
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}

	staffs := s.service.All()
	for i := 0; i < len(data); i++ {
		for _, staff := range staffs {
			if data[i].StaffUID == staff.Uid {
				data[i].StaffName = staff.GetCommonName()
			}
		}
	}
	apiOk(c, data, 0)
}

func (s *server) teamAdd(c *gin.Context) {
	var param teamAddParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	team := &team.Team{
		Name: param.Name,
	}
	if err := s.service.Team().Store(team); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, true, 0)
}

func (s *server) teamDelete(c *gin.Context) {
	var param teamDeleteParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	if err := s.service.Team().Delete(param.ID); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, true, 0)
}

func (s *server) teamMemberOp(c *gin.Context) {
	var param team.TeamOpParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	switch param.Op {
	case team.TeamOpAdd:
		if err := s.service.Team().AddMember(param.TeamID, param.UIDs...); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	case team.TeamOpRemove:
		if err := s.service.Team().RemoveMember(param.TeamID, param.UIDs...); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	default:
		apiError(c, ERROR_PARAM, "unknown operate")
		return
	}
	apiOk(c, true, 0)
}

func (s *server) teamManagerOp(c *gin.Context) {
	var param team.TeamOpParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	switch param.Op {
	case team.TeamOpAdd:
		if err := s.service.Team().AddManager(param.TeamID, param.UIDs[0]); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	case team.TeamOpRemove:
		if err := s.service.Team().RemoveManager(param.TeamID, param.UIDs[0]); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	default:
		apiError(c, ERROR_PARAM, "unknown operate")
		return
	}
	apiOk(c, true, 0)
}
