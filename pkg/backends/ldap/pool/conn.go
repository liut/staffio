package pool

import (
	"sync/atomic"
	"time"

	"github.com/go-ldap/ldap"
)

// Conn ...
type Conn struct {
	ldap.Client

	pooled    bool
	createdAt time.Time
	usedAt    atomic.Value
}

// NewConn ...
func NewConn(client ldap.Client) *Conn {
	cn := &Conn{
		Client:    client,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

// UsedAt ...
func (cn *Conn) UsedAt() time.Time {
	return cn.usedAt.Load().(time.Time)
}

// SetUsedAt ...
func (cn *Conn) SetUsedAt(tm time.Time) {
	cn.usedAt.Store(tm)
}

// // Close ...
// func (cn *Conn) Close() error {
// 	cn.Client.Close()
// 	return nil
// }
