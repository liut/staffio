package web

import (
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
	ID      int    `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type weeklyReportUpParam struct {
	ID int `json:"id" binding:"required"`
}

type weeklyReportStatusParam struct {
	UID    string        `json:"uid" binding:"required"`
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

	if err := s.service.Weekly().Add(user.UID, param.Content); err != nil {
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
	obj, err := s.service.Weekly().Get(param.ID)
	if err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	if obj.Uid != user.UID {
		apiError(c, ERROR_PARAM, "not yours")
		return
	}

	if err := s.service.Weekly().Update(param.ID, param.Content); err != nil {
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
	if err := s.service.Weekly().Applaud(param.ID, user.UID); err != nil {
		apiError(c, ERROR_DB, err)
		return
	}
	apiOk(c, nil, 0)
}

func (s *server) weeklyReportList(c *gin.Context) {
	var param weekly.ReportsSpec
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
		staff := staffs.WithUid(ret[i].Uid)
		if staff != nil {
			ret[i].Name = staff.GetCommonName()
			ret[i].Avatar = staff.AvatarUri()
		}
	}
	apiOk(c, ret, total)
}

func (s *server) weeklyReportListSelf(c *gin.Context) {
	var param weekly.ReportsSpec
	if err := c.Bind(&param); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	user := UserWithContext(c)
	param.UID = user.UID

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
	all := s.allStaffs(false)
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
	if len(param.Weeks) == 0 && param.Week > 0 {
		param.Weeks = []int{param.Week}
	}
	if err := s.service.Weekly().AddStatus(param.UID, status, param.Year, param.Weeks...); err != nil {
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
	if err := s.service.Weekly().RemoveStatus(param.ID); err != nil {
		apiError(c, ERROR_PARAM, err)
		return
	}
	apiOk(c, nil, 0)
}

type simpStaff struct {
	ID     int    `json:"id,omitempty"`
	Uid    string `json:"uid"`
	Name   string `json:"name"`
	Email  string `json:"email,omitempty"`
	Mobile string `json:"mobile,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	Created *time.Time `json:"created,omitempty"`
}

func (s *server) staffList(c *gin.Context) {
	ret := s.allStaffs(c.Request.FormValue("simple") != "yes")

	apiOk(c, ret, len(ret))
}

func (s *server) allStaffs(isFull bool) []*simpStaff {

	staffs := s.service.All()
	models.ByUid.Sort(staffs)
	var ret = make([]*simpStaff, len(staffs))
	for i, v := range staffs {
		ss := &simpStaff{
			Uid:     v.Uid,
			Name:    v.GetCommonName(),
			Created: v.Created,
		}
		if isFull {
			ss.ID = v.EmployeeNumber
			ss.Email = v.Email
			ss.Mobile = v.Mobile
			ss.Avatar = v.AvatarUri()
		}
		ret[i] = ss
	}
	return ret
}
