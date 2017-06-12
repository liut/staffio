package models

import (
	// "fmt"
	"sort"
	"strings"
)

var (
	ByUid = By(func(p1, p2 *Staff) bool {
		return p1.Uid < p2.Uid
	})
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

// employment for a person
type Staff struct {
	Uid            string `json:"uid" form:"uid" binding:"required"`         // 登录名
	Passwd         string `json:"-" form:"password"`                         // 密码
	CommonName     string `json:"cn,omitempty" form:"cn" binding:"required"` // 姓名（全名）
	GivenName      string `json:"gn" form:"gn" binding:"required"`           // 名 FirstName
	Surname        string `json:"sn" form:"sn" binding:"required"`           // 姓 LastName
	Nickname       string `json:"nickname,omitempty" form:"nickname"`        // 昵称
	Birthday       string `json:"birthday,omitempty" form:"birthday"`        // 生日
	Gender         Gender `json:"gender,omitempty" form:"gender"`            // 性别
	Email          string `json:"email" form:"email" binding:"required"`     // 邮箱
	Mobile         string `json:"mobile" form:"mobile" binding:"required"`   // 手机
	Tel            string `json:"tel,omitempty" form:"tel"`                  // 座机
	EmployeeNumber string `json:"eid,omitempty" form:"eid"`                  // 员工编号
	EmployeeType   string `json:"etype,omitempty" form:"etitle"`             // 员工岗位
	AvatarPath     string `json:"avatarPath,omitempty" form:"avatar"`        // 头像
	Description    string `json:"description,omitempty" form:"description"`  // 描述
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

// By is the type of a "less" function that defines the ordering of its Staff arguments.
type By func(p1, p2 *Staff) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(staffs []*Staff) {
	ps := &staffSorter{
		staffs: staffs,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type staffSorter struct {
	staffs []*Staff
	by     func(p1, p2 *Staff) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *staffSorter) Len() int {
	return len(s.staffs)
}

// Swap is part of sort.Interface.
func (s *staffSorter) Swap(i, j int) {
	s.staffs[i], s.staffs[j] = s.staffs[j], s.staffs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *staffSorter) Less(i, j int) bool {
	return s.by(s.staffs[i], s.staffs[j])
}
