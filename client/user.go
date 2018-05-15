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
	gash := time.Now().Unix() - u.LastHit
	if UserLifetime == 0 || gash > UserLifetime { // expired
		return false
	}

	return gash < UserLifetime && gash > UserLifetime/2
}

func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}
