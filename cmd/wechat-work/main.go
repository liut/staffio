// List or sync teams data from department of wechat work
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

var (
	action string
	uid    string

	nameReplacer = strings.NewReplacer("公司", "", "总部", "")
)

func init() {
	flag.StringVar(&action, "act", "", "action: list | sync | sync-all")
	flag.StringVar(&uid, "uid", "", "query uid")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	settings.Parse()
	// backends.Prepare()
	if action == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("action: %s", action)

	wechat := exwechat.New(settings.WechatCorpID, settings.WechatContactSecret)
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

	department, err := wechat.ListDepartment(1)
	if err != nil {
		log.Print(err)
		return
	}
	// log.Printf("department: %v", data)
	svc := backends.NewService()
	for _, dept := range department {
		// fmt.Printf("%4d  %10s   %d   %d\n", dept.Id, dept.Name, dept.ParentId, dept.Order)
		users, err := wechat.ListUser(dept.Id, false)
		if err != nil {
			log.Print(err)
			return
		}
		name := dept.Name

		if parent, _err := exwechat.FilterDepartment(department, dept.ParentId); _err == nil {
			name = nameReplacer.Replace(parent.Name) + "-" + dept.Name
		}
		var members []string
		var staffs []*models.Staff
		for _, val := range users {
			if !val.IsActived() || !val.IsEnabled() {
				log.Printf("user %s status %s, enabled %v", val.Name, val.Status, val.Enabled)
				continue
			}
			members = append(members, val.UID)
			staff := userToStaff(&val)
			// fmt.Println(staff)
			staffs = append(staffs, staff)
			// fmt.Printf("%4s %10s  %v\n", val.UID, val.Name, val.DepartmentIds)
		}

		fmt.Printf("%4d  %10s   %v \n", dept.Id, name, members)
		if action[:4] == "sync" {
			if action == "sync-all" {
				for _, staff := range staffs {
					if err = svc.SaveStaff(staff); err != nil {
						fmt.Printf("save staff %v ERR %s\n", staff, err)
						return
					}
					fmt.Printf("save staff %s(%s) OK\n", staff.CommonName, staff.Uid)
				}
			}
			err = svc.Team().Store(dept.Id, name, "", members)
			if err == nil {
				fmt.Printf("saved department to team OK\n")
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
