package models

import (
	"strings"
)

// SplitName 这个其实没有用
func SplitName(cn string) (sn, gn string) {
	cn = strings.TrimSpace(cn)
	if pos := strings.LastIndexByte(cn, ' '); pos > 0 {
		return cn[pos+1:], cn[0:pos]
	}
	a := strings.Split(strings.Trim(cn, " "), "")
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
