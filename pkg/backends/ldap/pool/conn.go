package pool

import (
	"sync/atomic"
	"time"

	"github.com/go-ldap/ldap"
)

type Conn struct {
	ldap.Client

	pooled    bool
	createdAt time.Time
	usedAt    atomic.Value
}

func NewConn(client ldap.Client) *Conn {
	cn := &Conn{
		Client:    client,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) UsedAt() time.Time {
	return cn.usedAt.Load().(time.Time)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	cn.usedAt.Store(tm)
}

// func (cn *Conn) Close() error {
// 	cn.Client.Close()
// 	return nil
// }
