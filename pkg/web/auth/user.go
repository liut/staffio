package auth

import (
	"encoding/base64"
	"log"
	"strings"
	"time"
)

//go:generate msgp -io=false

var (
	UserLifetime int64 = 3600
	Guest              = &User{}
)

// User 在线用户
type User struct {
	UID        string   `json:"uid" msg:"u"`
	Name       string   `json:"name" msg:"n"`
	Privileges string   `json:"privileges,omitempty" msg:"p"`
	LastHit    int64    `json:"hit,omitempty" msg:"h"`
	TeamID     int      `json:"tid,omitempty" msg:"t"`
	Watchings  []string `json:"watching,omitempty" msg:"w"`
}

func (u *User) IsExpired() bool {
	if UserLifetime <= 0 {
		return false
	}
	return u.LastHit+int64(UserLifetime) < time.Now().Unix()
}

func (u *User) NeedRefresh() bool {
	gash := time.Now().Unix() - u.LastHit
	if UserLifetime == 0 || gash > UserLifetime { // expired
		return false
	}

	return gash < UserLifetime && gash > UserLifetime/2
}

// refresh lastHit to time Unix
func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}

func (u User) Encode() (s string, err error) {
	var b []byte
	b, err = u.MarshalMsg(nil)
	if err == nil {
		s = strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
	}
	return
}

func (u *User) Decode(s string) (err error) {
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}

	var b []byte
	b, err = base64.URLEncoding.DecodeString(s)
	if err != nil {
		log.Printf("decode token %q ERR %s", s, err)
		return
	}

	*u = User{}
	_, err = u.UnmarshalMsg(b)
	if err != nil {
		log.Printf("unmarshal(%d) to msgpack from %q ERR %s", len(b), s, err)
	}

	return
}

type staff interface {
	GetUID() string
	GetName() string
}

func FromStaff(staff staff) *User {
	return &User{
		UID:  staff.GetUID(),
		Name: staff.GetName(),
	}
}
