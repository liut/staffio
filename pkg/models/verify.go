package models

import (
	"encoding/binary"
	"hash/crc32"
	"time"

	"lcgc/platform/staffio/pkg/models/common"
	"lcgc/platform/staffio/pkg/models/random"
)

var (
	VerifyLifeSeconds = 86400
)

func HashCode(code string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(code)))
}

// 用户验证，如邮箱、手机等
type Verify struct {
	Id          int              `db:"id" json:"id"`
	Uid         string           `db:"uid" json:"uid"`
	Target      string           `db:"target" json:"target"`
	Type        common.AliasType `db:"type_id" json:"type"`
	CodeHash    int64            `db:"code_hash" json:"-"`
	LifeSeconds int              `db:"life_seconds" json:"life_seconds"`
	Created     time.Time        `db:"created" json:"created"`
	Updated     time.Time        `db:"updated" json:"updated"`

	Code string `db:"-" json:"-"`
}

func (uv *Verify) IsExpired() bool {
	return time.Now().Unix() > uv.Updated.Unix()+int64(uv.LifeSeconds)
}

func (uv *Verify) Match(code string) bool {
	return HashCode(code) == uv.CodeHash
}

func (uv *Verify) CodeHashBytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(uv.CodeHash))
	return b
}

func NewVerify(at common.AliasType, target, uid string) *Verify {
	code := random.GenCode()
	codeHash := HashCode(code)
	return &Verify{
		Uid:         uid,
		Type:        at,
		Target:      target,
		LifeSeconds: VerifyLifeSeconds,
		CodeHash:    codeHash,
		Created:     time.Now(),
	}
}
