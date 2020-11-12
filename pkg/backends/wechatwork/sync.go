package wechatwork

import (
	"fmt"
	"sort"
	"strings"

	"fhyx.online/tencent-api-go/wxwork"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
	"github.com/liut/staffio/pkg/settings"
)

type User = wxwork.User
type Department = wxwork.Department
type Departments = wxwork.Departments

// Syncer ...
type Syncer struct {
	DeptFn    func(dept *wxwork.Department, idx int)
	BulkFn    func(svc backends.Servicer, t *team.Team, staffs models.Staffs) error
	WithTeam  bool
	WithStaff bool
	Output    bool

	api *wxwork.API
}

// SyncDepartment ...
func SyncDepartment(action, uid string) {
	s := &Syncer{WithTeam: strings.HasPrefix(action, "sync"), WithStaff: action == "sync-all"}
	s.api = wxwork.NewAPI(settings.Current.WechatCorpID, settings.Current.WechatContactSecret)
	s.BulkFn = backends.StoreTeamAndStaffs

	if action == "query" {
		if len(uid) > 0 {
			user, err := s.api.GetUser(uid)
			if err != nil {
				logger().Infow("get user fail", "err", err)
				return
			}
			if s.Output {
				fmt.Println(user)
			}

			return
		}
		logger().Infow("empty uid")
		return
	}
	s.RunIt()

}

// RunIt ...
func (s *Syncer) RunIt() error {
	if s.api == nil {
		s.api = wxwork.NewAPI(settings.Current.WechatCorpID, settings.Current.WechatContactSecret)
	}

	departments, err := s.api.ListDepartment("1")
	if err != nil {
		logger().Infow("list department fail", "err", err)
		return err
	}
	sort.Sort(departments)
	// log.Printf("departments: %v", data)
	svc := backends.NewService()
	for i, dept := range departments {
		if s.Output {
			fmt.Printf("%4d %4d %14s 	%8d\n", dept.ID, dept.ParentID, dept.Name, dept.Order)
		}

		if s.DeptFn != nil {
			s.DeptFn(&dept, i)
		}
		team := DepartmentToTeam(&dept, departments)
		var staffs models.Staffs

		if s.WithStaff {
			lr := wxwork.ListReq{DeptID: fmt.Sprint(dept.ID)}
			res, err := s.api.ListUser(lr)
			if err != nil {
				logger().Infow("list user fail", "err", err)
				return err
			}
			for _, val := range res.Users() {
				if !val.IsActived() || !val.IsEnabled() {
					logger().Infow("user actived?", "name", val.Name, "status", val.Status, "enabled", val.Enabled)
					continue
				}
				team.Members = append(team.Members, val.UID)
				staff := UserToStaff(&val)
				if val.IsLeader == 1 {
					team.Leaders = append(team.Leaders, staff.UID)
				}
				// fmt.Println(staff)
				staffs = append(staffs, *staff)
				// fmt.Printf("%4s %10s  %v\n", val.UID, val.Name, val.DepartmentIds)
			}

		}

		// fmt.Printf("%2d:%2d  %10s   %v \n", team.ID, team.ParentID, team.Name, team.Members)
		if s.WithTeam && s.BulkFn != nil {
			err = s.BulkFn(svc, team, staffs)
			if err != nil {
				logger().Infow("call bulkFn fail", "err", err)
				return err
			}
		}

	}
	logger().Infow("syncer ran all done")
	return nil
}
