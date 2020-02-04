package wechatwork

import (
	"fmt"
	"log"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

// SyncDepartment ...
func SyncDepartment(action, uid string) {

	wechat := exwechat.New(settings.Current.WechatCorpID, settings.Current.WechatContactSecret)
	if action == "query" {
		if len(uid) > 0 {
			user, err := wechat.GetUser(uid)
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

	departments, err := wechat.ListDepartment(1)
	if err != nil {
		log.Print(err)
		return
	}
	// log.Printf("departments: %v", data)
	svc := backends.NewService()
	for _, dept := range departments {
		// fmt.Printf("%4d  %10s   %d   %d\n", dept.Id, dept.Name, dept.ParentId, dept.Order)
		users, err := wechat.ListUser(dept.Id, false)
		if err != nil {
			log.Print(err)
			return
		}
		team := DepartmentToTeam(&dept, departments)

		var staffs []*models.Staff
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
			staffs = append(staffs, staff)
			// fmt.Printf("%4s %10s  %v\n", val.UID, val.Name, val.DepartmentIds)
		}

		fmt.Printf("%4d  %10s   %v \n", team.ID, team.Name, team.Members)
		if action[:4] == "sync" {
			var leader string
			if action == "sync-all" {
				for _, staff := range staffs {
					if err = svc.SaveStaff(staff); err != nil {
						fmt.Printf("save staff %v ERR %s\n", staff, err)
						return
					}
					fmt.Printf("save staff %s(%s) %q OK\n", staff.CommonName, staff.UID, leader)
				}
			}
			err = svc.Team().Store(team)
			if err == nil {
				fmt.Printf("saved team(%q, %q) to team OK\n", team.Name, team.Leaders)
				for _, leader := range team.Leaders {
					err = svc.Team().AddManager(team.ID, leader)
					if err != nil {
						log.Printf("add manager %q@%d ERR %s", leader, team.ID, err)
					}
				}
			} else {
				fmt.Printf("save team %s ERR %s\n", team.Name, err)
				return
			}
		}

	}
}
