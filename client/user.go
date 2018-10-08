//go:generate msgp

package client

import (
	"encoding/base64"
	"log"
	"strings"
	"time"
)

var (
	UserLifetime int64 = 3600
	Guest              = &User{}
)

type User struct {
	Uid     string `json:"uid" msg:"u"`
	Name    string `json:"name" msg:"n"`
	LastHit int64  `json:"hit,omitempty" msg:"h"`
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
