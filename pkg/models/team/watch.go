package team

import (
	"time"
)

// WatchStore interface of watch 关注
type WatchStore interface {
	Gets(uid string) Butts
	Watch(uid, target string) error
	Unwatch(uid, target string) error
}

// Butt 关注的对象
type Butt struct {
	UID    string `json:"uid" db:"watching"`
	Name   string `json:"name" db:"-"`
	Avatar string `json:"avatar" db:"-"`

	Created *time.Time `json:"created,omitempty" db:"created"`
}

// Butts ..
type Butts []Butt

// UIDs ...
func (z Butts) UIDs() []string {
	var arr = make([]string, len(z))
	for i, b := range z {
		arr[i] = b.UID
	}
	return arr
}
