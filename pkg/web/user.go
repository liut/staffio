package web

import (
	"time"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/models"
	. "lcgc/platform/staffio/pkg/settings"
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	LastHit int64  `json:"-"`
}

func (u *User) IsKeeper() bool {
	if u == nil {
		return false
	}
	return IsKeeper(u.Uid)
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

func (u *User) InGroup(gn string) bool {
	return InGroup(gn, u.Uid)
}

func IsKeeper(uid string) bool {
	return InGroup("keeper", uid)
}

func InGroup(gn, uid string) bool {
	return backends.InGroup(gn, uid)
}

func UserFromStaff(staff *models.Staff) *User {
	return &User{
		Uid:  staff.Uid,
		Name: staff.Name(),
	}
}
