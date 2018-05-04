package client

import (
	"time"
)

var (
	UserLifetime int64 = 3600
	Guest              = &User{}
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	LastHit int64  `json:"hit"`
}

func (u *User) IsExpired() bool {
	if UserLifetime == 0 {
		return false
	}
	return u.LastHit+UserLifetime < time.Now().Unix()
}

func (u *User) NeedRefresh() bool {
	return time.Now().Unix()-u.LastHit < UserLifetime/2
}

func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}
