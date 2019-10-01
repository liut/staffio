// List or sync teams data from department of wechat work
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
	// "github.com/liut/staffio/pkg/models/weekly"
	"github.com/liut/staffio/pkg/settings"
)

var (
	action string
	uid    string

	nameReplacer = strings.NewReplacer("公司", "", "总部", "")
)

func init() {
	flag.StringVar(&action, "act", "", "action: list | query | sync | sync-all")
	flag.StringVar(&uid, "uid", "", "query uid")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	if action == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("action: %q", action)

	// backends.InitSMTP()

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
		name := dept.Name

		if parent := departments.WithID(dept.ParentId); parent != nil {
			name = nameReplacer.Replace(parent.Name) + "-" + dept.Name
		}
		var team = &team.Team{
			ID:      dept.Id,
			Name:    name,
			Updated: time.Now(),
		}
		var staffs []*models.Staff
		for _, val := range users {
			if !val.IsActived() || !val.IsEnabled() {
				log.Printf("user %s status %s, enabled %v", val.Name, val.Status, val.Enabled)
				continue
			}
			team.Members = append(team.Members, val.UID)
			staff := userToStaff(&val)
			if val.IsLeader == 1 {
				team.Leaders = append(team.Leaders, staff.Uid)
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
					if staff.Leader {
						err = svc.Team().AddManager(dept.Id, staff.Uid)
						if err != nil {
							log.Printf("add manager %q@%d ERR %s", staff.Uid, dept.Id, err)
						}
						leader = staff.Uid
					}
					if err = svc.SaveStaff(staff); err != nil {
						fmt.Printf("save staff %v ERR %s\n", staff, err)
						return
					}
					fmt.Printf("save staff %s(%s) %q OK\n", staff.CommonName, staff.Uid, leader)
				}
			}
			err = svc.Team().Store(team)
			if err == nil {
				fmt.Printf("saved department(%q, %q) to team OK\n", team.Name, team.Leaders)
			} else {
				fmt.Printf("save department %s ERR %s\n", dept.Name, err)
				return
			}
		}

	}

}

func userToStaff(user *exwechat.User) *models.Staff {
	staff := &models.Staff{
		Uid:          strings.ToLower(user.UID),
		CommonName:   user.Name,
		Email:        user.Email,
		Mobile:       user.Mobile,
		Gender:       models.Gender(user.Gender),
		EmployeeType: user.Title,
		Leader:       user.IsLeader == 1,
	}
	staff.Surname, staff.GivenName = models.SplitName(user.Name)
	if user.Avatar != "" {
		uri, err := url.Parse(user.Avatar)
		if err == nil {
			staff.AvatarPath = uri.Path
		}
	}
	if user.Alias != "" {
		staff.Nickname = user.Alias
	}

	return staff
}
