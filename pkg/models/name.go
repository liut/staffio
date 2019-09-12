package models

import (
	"strings"
)

// SplitName 分隔姓名，使用 "·" " "
func SplitName(cn string) (sn, gn string) {
	cn = strings.TrimSpace(cn)
	if pos := strings.LastIndexByte(cn, ' '); pos > 0 {
		return cn[pos+1:], cn[0:pos]
	}
	if a := strings.Split(cn, "·"); len(a) == 2 {
		return a[0], a[1]
	}
	a := strings.Split(cn, "")
	switch len(a) {
	case 2:
		return a[0], a[1]
	case 3:
		return a[0], a[1] + a[2]
	case 4:
		return a[0] + a[1], a[2] + a[3]
	}
	return cn, ""
}
