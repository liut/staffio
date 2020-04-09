package team

import (
	"time"

	"github.com/liut/staffio/pkg/models/types"
)

type StringSlice = types.StringSlice

// RoleType ...
type RoleType int

const (
	RoleNothing RoleType = iota
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
	TeamID int        `json:"team_id" binding:"required" valid:"required"`
	UIDs   []string   `json:"staff_uids" binding:"required" valid:"required"`
}

// Store interface of team storage
type Store interface {
	// Get 取一个
	Get(id int) (*Team, error)
	// All 查询全部数据
	All(role RoleType) (data Teams, err error)
	// Store 保存
	Store(t *Team) error
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
	// GetWithMember
	GetWithMember(uid string) (*Team, error)
}

// Team work group
type Team struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	OrigName  string      `json:"orgname,omitempty" db:"-"`
	EngName   string      `json:"engname,omitempty" db:"-"`
	ParentID  int         `json:"parentID,omitempty" db:"parent_id"`
	Leaders   StringSlice `json:"leaders,omitempty"`
	Members   StringSlice `json:"members"`
	Created   time.Time   `json:"created,omitempty" db:"created"`
	Updated   *time.Time  `json:"updated,omitempty" db:"updated,omitempty"`
	StaffUID  string      `json:"staff_uid,omitempty" db:"staff_uid"`
	StaffName string      `json:"staff_name,omitempty" db:"-"`
	OrderNo   int         `json:"orderNo,omitempty" db:"-"`
}

// Teams ...
type Teams []Team

// UIDs return slice of staff uids
func (z Teams) UIDs() []string {
	var arr = make([]string, len(z))
	for i, t := range z {
		arr[i] = t.StaffUID
	}
	return arr
}
