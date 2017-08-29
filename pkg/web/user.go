package web

import (
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/liut/staffio/pkg/models"
)

type User struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	LastHit int64  `json:"-"`
	ver     uint8
}

const ver uint8 = 0

var (
	_ msgpack.CustomEncoder = (*User)(nil)
	_ msgpack.CustomDecoder = (*User)(nil)
)

func (u *User) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(ver, u.Uid, u.Name, u.LastHit)
}

func (u *User) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.Decode(&u.ver, &u.Uid, &u.Name, &u.LastHit)
}

func (u User) Encode() (b []byte, err error) {
	b, err = msgpack.Marshal(&u)
	return
}

func (u *User) Decode(b []byte) (err error) {
	err = msgpack.Unmarshal(b, u)
	return
}

func (u *User) IsExpired(lifetime int) bool {
	if lifetime == 0 {
		return false
	}
	return u.LastHit+int64(lifetime) < time.Now().Unix()
}

// refresh lastHit to time Unix
func (u *User) Refresh() {
	u.LastHit = time.Now().Unix()
}

func (u *User) IsKeeper() bool {
	for _, n := range keepers {
		if n == u.Uid {
			return true
		}
	}
	return false
}

func (s *server) IsKeeper(uid string) bool {
	return s.InGroup("keeper", uid)
}

func (s *server) InGroup(gn, uid string) bool {
	return s.service.InGroup(gn, uid)
}

func UserFromStaff(staff *models.Staff) *User {
	return &User{
		Uid:  staff.Uid,
		Name: staff.Name(),
	}
}
