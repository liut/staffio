package weekly

import (
	"encoding/json"
	"time"
)

type TeamRoleType int

const (
	RoleNothing TeamRoleType = iota
	RoleMember
	RoleManager
)

type TeamOpType int

const (
	TeamOpAdd    TeamOpType = 1 + iota // add
	TeamOpRemove                       // remove
)

type TeamOpParam struct {
	Op     TeamOpType `json:"op" binding:"required" valid:"[1:2]"`
	TeamID int        `json:"group_id" binding:"required" valid:"required"`
	UIDs   []string   `json:"staff_uids" binding:"required" valid:"required"`
}

// TeamStore interface of team storage
type TeamStore interface {
	// Get 取一个
	Get(id int) (*Team, error)
	// All 查询全部数据
	All(role TeamRoleType) (data []*Team, err error)
	// Store 保存
	Store(id int, name, leader string, members []string) error
	// Add members
	AddMember(id int, uids ...string) error
	// Remove members
	RemoveMember(id int, uids ...string) error
	// Delete 删除 Team
	Delete(id int) error
	// Add Manager
	AddManager(id int, uid string) error
	// Remove Manager
	RemoveManager(id int, uid string) error
}

// Team work group
type Team struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Leader    string          `json:"leader"`
	Members   json.RawMessage `json:"members"`
	Created   time.Time       `json:"created,omitempty" db:"created"`
	Updated   time.Time       `json:"updated,omitempty" db:"updated,omitempty"`
	StaffUID  string          `json:"staff_uid,omitempty" db:"staff_uid"`
	StaffName string          `json:"staff_name,omitempty" db:"-"`
}
