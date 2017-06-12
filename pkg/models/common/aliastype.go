package common

type AliasType uint8

const (
	AtEmail AliasType = 1 << iota // 1 邮箱
	AtPhone                       // 2 手机号
)

func (this AliasType) String() string {
	switch this {
	case AtEmail:
		return "email"
	case AtPhone:
		return "phone"
	}
	return ""
}
