package models

type Staff struct {
	Uid            string
	Passwd         string
	CommonName     string // 全名
	GivenName      string // 名
	SurName        string // 姓
	DisplayName    string // 昵称
	Email          string
	Mobile         string
	EmployeeNumber string
	EmployeeType   string
	Description    string
}

func (u *Staff) Name() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}

	if u.CommonName != "" {
		return u.CommonName
	}

	if u.SurName != "" && u.GivenName != "" {
		return u.SurName + u.GivenName
	}

	return u.Uid
}
