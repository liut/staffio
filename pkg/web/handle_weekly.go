package web

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/weekly"
)

const (
	ERROR_DB       = 1
	ERROR_PARAM    = 2
	ERROR_INTERNAL = 3
	ERROR_LIMIT    = 4
)

type weeklyReportAddParam struct {
	Content string `json:"content" binding:"required"`
}

type weeklyReportUpdateParam struct {
	Id      int    `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type weeklyReportUpParam struct {
	Id int `json:"id" binding:"required"`
}

type weeklyReportStatusParam struct {
	Uid    string        `json:"uid" binding:"required"`
	Year   int           `json:"year" `
	Week   int           `json:"week" `
	Weeks  []int         `json:"weeks" `
	Status weekly.Status `json:"status" `
}

func (s *server) weeklyReportAdd(c *gin.Context) {
	var param weeklyReportAddParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	user := UserWithContext(c)

	if err := s.service.Weekly().Add(user.Uid, param.Content); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) weeklyReportUpdate(c *gin.Context) {
	var param weeklyReportUpdateParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	user := UserWithContext(c)
	obj, err := s.service.Weekly().Get(param.Id)
	if err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	if obj.Uid != user.Uid {
		apiError(c, ERROR_PARAM, "not yours")
		return
	}

	if err := s.service.Weekly().Update(param.Id, param.Content); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) weeklyReportUp(c *gin.Context) {
	var param weeklyReportUpParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	user := UserWithContext(c)
	if err := s.service.Weekly().Applaud(param.Id, user.Uid); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) weeklyReportList(c *gin.Context) {
	var param weekly.ListParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}

	ret, total, err := s.service.Weekly().All(param)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	staffs := s.service.All()
	for i := 0; i < len(ret); i++ {
		for _, staff := range staffs {
			if staff.Uid == ret[i].Uid {
				ret[i].Name = staff.GetCommonName()
			}
		}

	}
	apiOk(c, ret, total)
}

func (s *server) weeklyReportListSelf(c *gin.Context) {
	var param weekly.ListParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	user := UserWithContext(c)
	param.Uid = user.Uid

	ret, total, err := s.service.Weekly().All(param)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, ret, total)
}

func (s *server) weeklyReportStat(c *gin.Context) {
	var param weekly.ReportStatParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	start, err1 := formatDate(param.Start, "00:00:01")
	end, err2 := formatDate(param.End, "23:59:59")

	if err1 != nil || err2 != nil {
		errMsg := ""
		if err1 != nil {
			errMsg += err1.Error()
		}
		if err2 != nil {
			errMsg += " " + err2.Error()
		}
		apiError(c, ERROR_PARAM, errMsg)
		return
	}
	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = time.Date(end.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}
	start = getWeekFirstDate(start)
	end = getWeekLastDate(end)
	ret, err := s.service.Weekly().Stat(start, end)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	all := s.allStaffs()
	ignores, err := s.service.Weekly().StatusRecords(weekly.WRIgnore)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	for _, staff := range all {
		var isIgnore bool
		for _, ig := range ignores {
			if ig.Uid == staff.Uid {
				// 		ig.Name = staff.Name
				isIgnore = true
				continue
			}
		}
		if isIgnore {
			continue
		}
		ret.All = append(ret.All, &weekly.ReportUser{
			Uid:     staff.Uid,
			Name:    staff.Name,
			Created: staff.Created,
		})
	}
	apiOk(c, ret, 1)
}

func formatDate(dateStr, timeStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
}

func getWeekFirstDate(date time.Time) time.Time {
	day := int(date.Weekday())
	year, _ := date.ISOWeek()
	var start time.Time
	if year >= date.Year() {
		if day == 0 {
			start = date.AddDate(0, 0, -6)
		} else {
			start = date.AddDate(0, 0, 1-day)
		}
	} else {
		if day == 0 {
			start = date.AddDate(0, 0, 1)
		} else {
			start = date.AddDate(0, 0, 8-day)
		}
	}
	return start
}

func getWeekLastDate(date time.Time) time.Time {
	day := int(date.Weekday())
	year, _ := date.ISOWeek()
	var end time.Time
	if year <= date.Year() {
		if day == 0 {
			end = date
		} else {
			end = date.AddDate(0, 0, 7-day)
		}
	} else {
		if day == 0 {
			end = date.AddDate(0, 0, -7)
		} else {
			end = date.AddDate(0, 0, -day)
		}
	}
	if time.Now().Before(end) {
		end = time.Now()
	}
	return end
}

func (s *server) weeklyProblemList(c *gin.Context) {
	// TODO: list problem
	apiOk(c, nil, 0)
}

func (s *server) weeklyProblemAdd(c *gin.Context) {
	// TODO: add problem
	apiOk(c, nil, 0)
}

func (s *server) weeklyProblemUpdate(c *gin.Context) {
	// TODO: update problem
	apiOk(c, nil, 0)
}

func (s *server) weeklyVacationList(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		apiError(c, ERROR_PARAM, "empty uid")
		return
	}
	data, err := s.service.Weekly().StatusRecordsWithUser(weekly.WRVacation, uid)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, data, len(data))
}

func (s *server) weeklyVacationAdd(c *gin.Context) { s.weeklyStatusAdd(c, weekly.WRVacation) }

func (s *server) weeklyVacationRemove(c *gin.Context) { s.weeklyStatusRemove(c) }

func (s *server) weeklyIgnoreList(c *gin.Context) {
	data, err := s.service.Weekly().StatusRecords(weekly.WRIgnore)
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, data, len(data))
}

func (s *server) weeklyIgnoreAdd(c *gin.Context) { s.weeklyStatusAdd(c, weekly.WRIgnore) }

func (s *server) weeklyIgnoreRemove(c *gin.Context) { s.weeklyStatusRemove(c) }

func (s *server) weeklyStatusAdd(c *gin.Context, status weekly.Status) {
	var param weeklyReportStatusParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	if len(param.Weeks) > 0 && param.Week == 0 {
		param.Week = param.Weeks[0]
	}
	if err := s.service.Weekly().AddStatus(param.Uid, param.Year, param.Week, status); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) weeklyStatusRemove(c *gin.Context) {
	var param idParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	if err := s.service.Weekly().RemoveStatus(param.Id); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) teamListByRole(c *gin.Context) {
	role, err := strconv.Atoi(c.Request.FormValue("role"))
	if err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	data, err := s.service.Team().All(weekly.TeamRoleType(role))
	if err != nil {
		apiError(c, ERROR_DB, err)
		return
	}

	staffs := s.service.All()
	for i := 0; i < len(data); i++ {
		for _, staff := range staffs {
			if data[i].StaffUid == staff.Uid {
				data[i].StaffName = staff.GetCommonName()
			}
		}
	}
	apiOk(c, data, 0)
}

func (s *server) teamMemberOp(c *gin.Context) {
	var param weekly.TeamOpParam
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	switch param.Op {
	case weekly.TeamOpAdd:
		if err := s.service.Team().AddMember(param.TeamId, param.Uids...); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	case weekly.TeamOpRemove:
		if err := s.service.Team().RemoveMember(param.TeamId, param.Uids...); err != nil {
			apiError(c, ERROR_DB, err)
			return
		}
	default:
		apiError(c, ERROR_PARAM, "unknown operate")
		return
	}
	apiOk(c, true, 0)
}

type simpStaff struct {
	Id      int       `json:"id"`
	Uid     string    `json:"uid"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Mobile  string    `json:"mobile"`
	Created time.Time `json:"created,omitempty"`
}

func (s *server) staffList(c *gin.Context) {
	ret := s.allStaffs()
	apiOk(c, ret, len(ret))
}

func (s *server) allStaffs() []*simpStaff {

	staffs := s.service.All()
	models.ByUid.Sort(staffs)
	var ret = make([]*simpStaff, len(staffs))
	for i, v := range staffs {
		ret[i] = &simpStaff{
			Id:      v.EmployeeNumber,
			Uid:     v.Uid,
			Name:    v.GetCommonName(),
			Email:   v.Email,
			Mobile:  v.Mobile,
			Created: v.Created,
		}
	}
	return ret
}
