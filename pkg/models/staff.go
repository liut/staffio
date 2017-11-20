package models

import (
	"strings"
)

var (
	ProfileEditables = map[string]string{
		"nickname":    "displayName",
		"cn":          "cn",
		"gn":          "givenName",
		"sn":          "sn",
		"email":       "mail",
		"mobile":      "mobile",
		"eid":         "employeeNumber",
		"etitle":      "employeeType",
		"birthday":    "dateOfBirth",
		"gender":      "gender",
		"avatarPath":  "avatarPath",
		"description": "description",
	}

	cnFormat = "<gn> <sn>"
)

func SetNameFormat(s string) {
	cnFormat = s
}

// employment for a person
type Staff struct {
	Uid            string `json:"uid" form:"uid" binding:"required"`        // 登录名
	CommonName     string `json:"cn,omitempty" form:"cn"`                   // 姓名（全名）
	GivenName      string `json:"gn" form:"gn" binding:"required"`          // 名 FirstName
	Surname        string `json:"sn" form:"sn" binding:"required"`          // 姓 LastName
	Nickname       string `json:"nickname,omitempty" form:"nickname"`       // 昵称
	Birthday       string `json:"birthday,omitempty" form:"birthday"`       // 生日
	Gender         Gender `json:"gender,omitempty" form:"gender"`           // 性别
	Email          string `json:"email" form:"email" binding:"required"`    // 邮箱
	Mobile         string `json:"mobile" form:"mobile" binding:"required"`  // 手机
	Tel            string `json:"tel,omitempty" form:"tel"`                 // 座机
	EmployeeNumber string `json:"eid,omitempty" form:"eid"`                 // 员工编号
	EmployeeType   string `json:"etype,omitempty" form:"etitle"`            // 员工岗位
	AvatarPath     string `json:"avatarPath,omitempty" form:"avatar"`       // 头像
	Description    string `json:"description,omitempty" form:"description"` // 描述
}

func (u *Staff) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}

	if u.CommonName != "" {
		return u.CommonName
	}

	if u.Surname != "" && u.GivenName != "" {
		return formatCN(u.GivenName, u.Surname)
	}

	return u.Uid
}

func (u *Staff) GetCommonName() string {
	if u.CommonName != "" {
		return u.CommonName
	}

	return formatCN(u.GivenName, u.Surname)
}

// func (u *Staff) String() string {
// 	name := u.Name()
// 	if name == u.Uid {
// 		return name
// 	}

// 	return fmt.Sprintf("%s (%s)", name, u.Uid)
// }

func formatCN(gn, sn string) string {
	r := strings.NewReplacer("<gn>", gn, "<sn>", sn)
	return r.Replace(cnFormat)
}
