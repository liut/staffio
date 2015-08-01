package web

import (
	"time"
	"tuluu.com/liut/staffio/backends"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	LastHit int64  `json:"last_hit"`
}

func (u *User) IsKeeper() bool {
	if u == nil {
		return false
	}
	return backends.InGroup("keeper", u.Uid)
}

func (u *User) IsExpired() bool {
	lifetime := Settings.UserLifetime
	if lifetime == 0 {
		return false
	}
	return u.LastHit+int64(lifetime) < time.Now().Unix()
}

// refresh lastHit to time Unix
func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}

func UserFromStaff(staff *models.Staff) *User {
	return &User{
		Uid:  staff.Uid,
		Name: staff.Name(),
	}
}
