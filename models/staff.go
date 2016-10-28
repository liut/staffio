package models

import (
	// "fmt"
	"sort"
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
)

func NewStaff(uid, cn, email string) *Staff {
	sn, gn := SplitName(cn)
	return &Staff{
		Uid:        uid,
		CommonName: cn,
		Surname:    sn,
		GivenName:  gn,
		Email:      email,
	}
}

type Staff struct {
	Uid            string `json:"uid" form:"uid" binding:"required"` // 登录名
	Passwd         string `json:"-" form:"password"`
	CommonName     string `json:"cn,omitempty"`                       // 全名
	GivenName      string `json:"gn" form:"gn" binding:"required"`    // 名
	Surname        string `json:"sn" form:"sn" binding:"required"`    // 姓
	Nickname       string `json:"nickname,omitempty" form:"nickname"` // 昵称
	Birthday       string `json:"birthday,omitempty" form:"birthday"`
	Gender         Gender `json:"gender,omitempty" form:"gender"`
	Email          string `json:"email" form:"email" binding:"required"`
	Mobile         string `json:"mobile" form:"mobile" binding:"required"`
	EmployeeNumber string `json:"eid,omitempty" form:"eid"`
	EmployeeType   string `json:"etype,omitempty" form:"etitle"`
	AvatarPath     string `json:"avatarPath,omitempty" form:"avatar"`
	Description    string `json:"description,omitempty" form:"description"`
}

func (u *Staff) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}

	if u.CommonName != "" {
		return u.CommonName
	}

	if u.Surname != "" && u.GivenName != "" {
		return u.Surname + u.GivenName
	}

	return u.Uid
}

func (u *Staff) GetCommonName() string {
	if u.CommonName != "" {
		return u.CommonName
	}

	return u.Surname + u.GivenName
}

// func (u *Staff) String() string {
// 	name := u.Name()
// 	if name == u.Uid {
// 		return name
// 	}

// 	return fmt.Sprintf("%s (%s)", name, u.Uid)
// }

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
