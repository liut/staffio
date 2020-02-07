package wechatwork

import (
	"net/url"
	"strings"

	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
)

var nameReplacer = strings.NewReplacer("公司", "", "总部", "", "分公司", "", "项目组", "")

// UserToStaff ...
func UserToStaff(user *exwechat.User) *models.Staff {
	staff := &models.Staff{
		UID:          strings.ToLower(user.UID),
		CommonName:   user.Name,
		Email:        user.Email,
		Mobile:       user.Mobile,
		Tel:          user.Tel,
		Gender:       models.Gender(user.Gender).String(),
		EmployeeType: user.Title,
		// Leader:       user.IsLeader == 1,
	}
	fullname := user.Name
	if user.EnglishName != "" {
		fullname = user.EnglishName
	}
	staff.Surname, staff.GivenName = models.SplitName(fullname)

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

// DepartmentToTeam ...
func DepartmentToTeam(dept *exwechat.Department, all exwechat.Departments) *team.Team {
	var team = &team.Team{
		ID:       dept.Id,
		Name:     dept.Name,
		OrigName: dept.Name,
		ParentID: dept.ParentId,
		OrderNo:  dept.Order,
	}

	if all != nil {
		if parent := all.WithID(dept.ParentId); parent != nil {
			team.Name = nameReplacer.Replace(parent.Name) + "-" + dept.Name
		}
	}

	return team
}
