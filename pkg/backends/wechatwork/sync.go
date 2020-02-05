package wechatwork

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
	"github.com/liut/staffio/pkg/settings"
)

// Syncer ...
type Syncer struct {
	BulkFn    func(svc backends.Servicer, t *team.Team, staffs models.Staffs) error
	WithTeam  bool
	WithStaff bool

	api *exwechat.API
}

// SyncDepartment ...
func SyncDepartment(action, uid string) {
	s := &Syncer{WithTeam: strings.HasPrefix(action, "sync"), WithStaff: action == "sync-all"}
	s.api = exwechat.New(settings.Current.WechatCorpID, settings.Current.WechatContactSecret)
	s.BulkFn = backends.StoreTeamAndStaffs

	if action == "query" {
		if len(uid) > 0 {
			user, err := s.api.GetUser(uid)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Println(user)
			return
		}
		fmt.Println("empty uid")
		return
	}
	s.RunIt()

}

// RunIt ...
func (s *Syncer) RunIt() error {
	if s.api == nil {
		s.api = exwechat.New(settings.Current.WechatCorpID, settings.Current.WechatContactSecret)
	}

	departments, err := s.api.ListDepartment(1)
	if err != nil {
		log.Print(err)
		return err
	}
	sort.Sort(departments)
	// log.Printf("departments: %v", data)
	svc := backends.NewService()
	for _, dept := range departments {
		fmt.Printf("%4d %4d %14s 	%8d\n", dept.Id, dept.ParentId, dept.Name, dept.Order)
		team := DepartmentToTeam(&dept, departments)
		var staffs models.Staffs

		if s.WithStaff {
			users, err := s.api.ListUser(dept.Id, false)
			if err != nil {
				log.Print(err)
				return err
			}
			for _, val := range users {
				if !val.IsActived() || !val.IsEnabled() {
					log.Printf("user %s status %s, enabled %v", val.Name, val.Status, val.Enabled)
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
				log.Print(err)
				return err
			}
		}

	}
	log.Print("all done")
	return nil
}
