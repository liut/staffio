package client

// Staff is a retrieved employee struct.
type Staff struct {
	Uid            string `json:"uid" form:"uid"`                     // 登录名
	CommonName     string `json:"cn,omitempty" form:"cn"`             // 全名
	GivenName      string `json:"gn" form:"gn"`                       // 名
	Surname        string `json:"sn" form:"sn"`                       // 姓
	Nickname       string `json:"nickname,omitempty" form:"nickname"` // 昵称
	Birthday       string `json:"birthday,omitempty" form:"birthday"` // 生日
	Gender         uint8  `json:"gender,omitempty"`                   // 1=male, 2=female, 0=unknown
	Mobile         string `json:"mobile,omitempty"`                   // cell phone number
	Email          string `json:"email,omitempty"`
	EmployeeNumber string `json:"eid,omitempty" form:"eid"`
	EmployeeType   string `json:"etype,omitempty" form:"etitle"`
	AvatarPath     string `json:"avatarPath,omitempty" form:"avatar"`
	Provider       string `json:"provider,omitempty"`
}

type RoleMe map[string]interface{}

func (r RoleMe) Has(name string) bool {
	if v, exist := r[name]; exist {
		if g, ok := v.(bool); ok {
			return g
		}
	}
	return false
}
