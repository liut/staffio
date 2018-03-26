package models

import (
	"bytes"
)

type Gender uint8

const (
	Unknown Gender = 0 + iota
	Male
	Female
)

var genderKeys = []string{"unknown", "male", "female"}

func (this Gender) String() string {
	if this >= Unknown && this <= Female {
		return genderKeys[this]
	}
	return "unknown"
}

func (this Gender) MarshalText() ([]byte, error) {
	switch this {
	case Male:
		return []byte{'m'}, nil
	case Female:
		return []byte{'f'}, nil
	}
	return []byte{'u'}, nil
}

func (this *Gender) UnmarshalText(b []byte) (err error) {
	if len(b) == 0 {
		*this = Unknown
		return
	}
	r := bytes.Runes(bytes.Trim(b, "\""))
	switch c := r[0]; c {
	case 'm', 'M', '1', '男':
		*this = Male
	case 'f', 'F', '2', '女':
		*this = Female
	case 'u', 'U', '0':
		*this = Unknown
	default:
		err = ErrInvalidGender
	}
	return
}
