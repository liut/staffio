package models

import (
	"strings"
)

func SplitName(cn string) (sn, gn string) {
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
