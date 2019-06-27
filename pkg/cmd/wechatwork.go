// Copyright © 2019 liut <liutao@liut.cc>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/weekly"
)

// wechatworkCmd represents the wechatwork command
var wechatworkCmd = &cobra.Command{
	Use:   "wechatwork",
	Short: "Sync with wechat work",
	Run:   wechatworkRun,
}
var nameReplacer = strings.NewReplacer("公司", "", "总部", "", "分公司", "", "项目组", "")

func init() {
	RootCmd.AddCommand(wechatworkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wechatworkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wechatworkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	wechatworkCmd.Flags().String("act", "", "action: query | sync | sync-all")
	wechatworkCmd.Flags().String("uid", "", "query uid")
}

func wechatworkRun(cmd *cobra.Command, args []string) {
	action, _ := cmd.Flags().GetString("act")
	uid, _ := cmd.Flags().GetString("uid")

	log := logger()

	wechat := exwechat.New(settings.WechatCorpID, settings.WechatContactSecret)
	if action == "query" {
		if len(uid) > 0 {
			user, err := wechat.GetUser(uid)
			if err != nil {
				log.Warnw("wechat GetUser fail", "err", err)
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
		log.Infow("wechat ListDepartment fail", "err", err)
		return
	}
	// log.Infof("departments: %v", data)
	svc := backends.NewService()
	for _, dept := range departments {
		// fmt.Printf("%4d  %10s   %d   %d\n", dept.Id, dept.Name, dept.ParentId, dept.Order)
		users, err := wechat.ListUser(dept.Id, false)
		if err != nil {
			log.Infow("wechat ListUser fail", "err", err)
			return
		}
		name := dept.Name

		if parent := departments.WithID(dept.ParentId); parent != nil {
			name = nameReplacer.Replace(parent.Name) + "-" + dept.Name
		}
		var team = &weekly.Team{
			ID:      dept.Id,
			Name:    name,
			Updated: time.Now(),
		}
		var staffs []*models.Staff
		for _, val := range users {
			if !val.IsActived() || !val.IsEnabled() {
				log.Infow("invalid user status", "user", val.Name, "status", val.Status, "enabled", val.Enabled)
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

		fmt.Printf("%4d  %12s   %v \n", team.ID, team.Name, team.Members)
		if strings.HasPrefix(action, "sync") {
			var leader string
			if action == "sync-all" {
				for _, staff := range staffs {
					if staff.Leader {
						err = svc.Team().AddManager(dept.Id, staff.Uid)
						if err != nil {
							log.Infow("add manager fail", "uid", staff.Uid, "id", dept.Id, "err", err)
						}
						leader = staff.Uid
					}
					if err = svc.SaveStaff(staff); err != nil {
						fmt.Printf("save staff %v ERR %s\n", staff, err)
						if err != backends.ErrEmptyEmail {
							return
						}
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
