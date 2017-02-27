package models

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

type AliasType uint8

const (
	AtEmail AliasType = 1 << iota // 1 邮箱
	AtPhone                       // 2 手机号
)

var (
	VerifyLifeSeconds = 86400
)

func GenVerifyCode() string {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	d := r.Intn(999999)
	return fmt.Sprintf("%06d", d)
}

func VerifyHashCode(code string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(code)))
}

// 用户验证，如邮箱、手机等
type Verify struct {
	Id          int       `db:"id" json:"id"`
	Uid         string    `db:"uid" json:"uid"`
	Target      string    `db:"target" json:"target"`
	Type        AliasType `db:"type_id" json:"type"`
	CodeHash    int64     `db:"code_hash" json:"-"`
	LifeSeconds int       `db:"life_seconds" json:"life_seconds"`
	Created     time.Time `db:"created" json:"created"`
	Updated     time.Time `db:"updated" json:"updated"`

	Code string `db:"-" json:"-"`
}

func (uv *Verify) IsExpired() bool {
	return time.Now().Unix() > uv.Updated.Unix()+int64(uv.LifeSeconds)
}

func (uv *Verify) Match(code string) bool {
	return VerifyHashCode(code) == uv.CodeHash
}

func (uv *Verify) CodeHashBytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(uv.CodeHash))
	return b
}

func NewVerify(at AliasType, target, uid string) *Verify {
	code := GenVerifyCode()
	codeHash := VerifyHashCode(code)
	return &Verify{
		Uid:         uid,
		Type:        at,
		Target:      target,
		LifeSeconds: VerifyLifeSeconds,
		CodeHash:    codeHash,
		Created:     time.Now(),
	}
}
