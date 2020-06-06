// Package sync command
package main

import (
	"flag"
	"log"

	"github.com/mozillazg/go-slugify"
	"go.uber.org/zap"

	"fhyx.online/welink-api-go/gender"
	wlog "fhyx.online/welink-api-go/log"
	"fhyx.online/welink-api-go/welink"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/backends/wechatwork"
	zlog "github.com/liut/staffio/pkg/log"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
	// "github.com/liut/staffio/pkg/settings"
)

var (
	version = "dev"
	action  string
)

func init() {
	flag.StringVar(&action, "act", "", "action: dept-list")
}

func logger() zlog.Logger {
	return zlog.GetLogger()
}

func inDevelop() bool {
	return version == "dev"
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	var logger *zap.Logger
	if inDevelop() {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	zlog.SetLogger(sugar)
	wlog.SetLogger(sugar)

	if action == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("action: %q", action)

	switch action {
	case "dept-sync":
		syncDepartment(false)
		break
	case "dept-status":
		syncDepartment(true)
		break
	}
}

const (
	corpSuffix = "@phoenixtv"
)

// type bulkItem struct {
// 	team   *team.Team
// 	staffs models.Staffs
// }

func syncDepartment(showStatus bool) {
	var departmentUps []welink.DepartmentUp
	var userUps []welink.UserUp
	s := &wechatwork.Syncer{
		WithTeam:  true,
		WithStaff: true,
	}

	s.BulkFn = func(svc backends.Servicer, team *team.Team, staffs models.Staffs) error {
		// logger().Infow("bulk ", "team", team)
		if team.ParentID == 0 {
			return nil
		}

		up := teamToWelinkDeptUp(team)
		if up.Level == 1 || up.Level == 2 || up.Level == 3 {
			departmentUps = append(departmentUps, *up)
		}
		for _, staff := range staffs {
			userUps = append(userUps, *staffToWelinkUserUp(&staff, team.ID))
		}
		return nil
	}

	err := s.RunIt()
	if err != nil {
		logger().Infow("run fail", "err", err)
		return
	}

	api := welink.NewAPI()
	var count int
	count = len(departmentUps)
	for i := 0; i < count; i += 10 {
		j := i + 10
		if j >= count {
			j = count - 1
		}
		ups := departmentUps[i:j]
		logger().Infow("data for depts up", "i", i, "j", j, "data", ups)
		var res []welink.DeptRespItem
		if showStatus {
			res, err = api.StatusDepartment(ups)
		} else {
			res, err = api.SyncDepartment(ups)
		}

		if err != nil {
			logger().Infow("sync fail", "len", len(ups), "err", err)
			return
		}
		logger().Infow("sync ok", "len", len(ups))
		for _, info := range res {
			logger().Infow("resp ", "info", info)
		}
	}

	count = len(userUps)
	for i := 0; i < count; i += 10 {
		j := i + 10
		if j >= count {
			j = count - 1
		}
		ups := userUps[i:j]
		logger().Infow("data for users up", "i", i, "j", j, "data", ups)
		var res []welink.UserRespItem
		if showStatus {
			res, err = api.StatusUser(ups)
		} else {
			res, err = api.SyncUser(ups)
		}

		if err != nil {
			logger().Infow("sync fail", "len", len(ups), "err", err)
			return
		}
		logger().Infow("sync ok", "len", len(ups))
		for _, info := range res {
			logger().Infow("resp ", "info", info)
		}
	}

}

func staffToWelinkUserUp(staff *models.Staff, teamID int) *welink.User {
	user := &welink.UserUp{
		CorpUID:            staff.UID,
		CorpDeptID:         teamID,
		NameCN:             staff.GetCommonName(),
		NameEN:             slugify.Slugify(staff.Name()),
		Gender:             gender.ParseGender([]byte(staff.Gender)),
		Mobile:             staff.Mobile,
		Phone:              staff.Mobile,
		Email:              staff.Email,
		IsOpenAccount:      1,
		IsHideMobileNumber: 2,
		Valid:              1,
		OrderInDepts:       1,
	}

	return user
}

func teamToWelinkDeptUp(team *team.Team) *welink.DepartmentUp {
	up := &welink.DepartmentUp{
		CorpDeptID:   team.ID,
		CorpParentID: team.ParentID,
		NameCN:       team.OrigName,
		NameEN:       slugify.Slugify(team.OrigName),
		Valid:        1,
		Level:        getTeamLevel(team.ParentID),
		OrderNo:      getOrderNo(team.OrderNo),
	}
	if team.ParentID == 1 {
		up.CorpParentID = 0
	}
	if len(team.Leaders) > 0 {
		up.Leader = team.Leaders[0] + corpSuffix
	}
	return up
}

func deptToWelinkDepartmentUp(dept *wechatwork.Department) *welink.DepartmentUp {
	up := &welink.DepartmentUp{
		CorpDeptID:   dept.Id,
		CorpParentID: dept.ParentId,
		NameCN:       dept.Name,
		NameEN:       slugify.Slugify(dept.Name),
		Valid:        1,
		Level:        getTeamLevel(dept.ParentId),
		OrderNo:      getOrderNo(dept.Order),
	}
	if dept.ParentId == 1 {
		up.CorpParentID = 0
	}

	return up
}

func getTeamLevel(parentID int) int {
	if parentID == 1 {
		return 1
	}
	if parentID == 6 || parentID == 7 || parentID == 23 {
		return 2
	}
	return 3
}

func getOrderNo(no int) int {
	return 1000 - int(float32(no)/float32(100005000)*999)
}
