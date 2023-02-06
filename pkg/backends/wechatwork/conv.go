package wechatwork

import (
	"net/url"
	"strings"

	"daxv.cn/gopak/tencent-api-go/wxwork"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
)

var nameReplacer = strings.NewReplacer("公司", "", "总部", "", "分公司", "", "项目组", "")

// UserToStaff ...
func UserToStaff(user *wxwork.User) *models.Staff {
	staff := &models.Staff{
		UID:          strings.ToLower(user.UID),
		CommonName:   user.Name,
		Mobile:       user.Mobile,
		Tel:          user.Tel,
		Gender:       models.Gender(user.Gender).String(),
		EmployeeType: user.Title,
		// Leader:       user.IsLeader == 1,
	}
	if len(user.BizMail) > 0 {
		staff.Email = user.BizMail
	} else {
		staff.Email = user.Email
	}
	fullname := user.Name
	staff.Surname, staff.GivenName = models.SplitName(fullname)

	if user.Avatar != "" {
		uri, err := url.Parse(user.Avatar)
		if err == nil {
			staff.AvatarPath = "//" + uri.Host + uri.Path
		}
	}
	if user.Alias != "" {
		staff.Nickname = user.Alias
	}

	return staff
}

// DepartmentToTeam ...
func DepartmentToTeam(dept *wxwork.Department, all wxwork.Departments) *team.Team {
	var team = &team.Team{
		ID:       int(dept.ID),
		Name:     dept.Name,
		OrigName: dept.Name,
		ParentID: int(dept.ParentID),
		OrderNo:  int(dept.Order),
	}

	if all != nil {
		if parent := all.WithID(int(dept.ParentID)); parent != nil {
			team.Name = nameReplacer.Replace(parent.Name) + "-" + dept.Name
		}
	}

	return team
}
