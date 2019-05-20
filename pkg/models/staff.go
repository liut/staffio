package models

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/liut/staffio/pkg/common"
)

type Gender = common.Gender

const (
	Unknown = common.Unknown
	Male    = common.Male
	Female  = common.Female
)

var (
	// ProfileEditables deprecated
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

	avatarReplacer = strings.NewReplacer("/0", "/60")
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
	EmployeeNumber int    `json:"eid,omitempty" form:"eid"`                 // 员工编号
	EmployeeType   string `json:"etype,omitempty" form:"etitle"`            // 员工岗位
	AvatarPath     string `json:"avatarPath,omitempty" form:"avatar"`       // 头像
	JpegPhoto      []byte `json:"-" form:"-"`                               // jpegPhoto data
	Description    string `json:"description,omitempty" form:"description"` // 描述
	JoinDate       string `json:"joinDate,omitempty" form:"joinDate"`       // 加入日期
	IDCN           string `json:"idcn,omitempty" form:"idcn"`               // 身份证号

	Created time.Time `json:"created,omitempty" form:"created"` // 创建时间
	Updated time.Time `json:"updated,omitempty" form:"updated"` // 修改时间

	DN string `json:"-" form:"-"` // distinguishedName of LDAP entry

	Leader bool `json:"leader,omitempty" form:"-"` // temporary var
	TeamID int  `json:"teamID,omitempty" form:"-"` // department id
}

func (u *Staff) GetUID() string {
	return u.Uid
}

func (u *Staff) GetName() string {
	return u.Name()
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

func (u *Staff) AvatarUri() string {
	if len(u.JpegPhoto) > 0 {
		return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(u.JpegPhoto)
	}
	if len(u.AvatarPath) > 0 {
		s := u.AvatarPath
		if strings.HasSuffix(s, "/") {
			s = s + "0"
		}
		return "https://p.qlogo.cn" + avatarReplacer.Replace(s)
	}
	return ""
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

type Staffs []*Staff

func (arr Staffs) WithUid(uid string) *Staff {
	for _, u := range arr {
		if u.Uid == uid {
			return u
		}
	}
	return nil
}
