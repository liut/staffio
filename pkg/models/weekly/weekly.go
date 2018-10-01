package weekly

import (
	"encoding/json"
	"time"
)

// 查询周报参数
type ListParam struct {
	Uid     string     `json:"uid"`
	GroupId int        `json:"group_id"`
	Pager   *ListPager `json:"pager" valid:"required"`
	Sort    *ListSort  `json:"sort" valid:"required"` // weekly_report.id,work_group_id
}

type WeeklyStore interface {
	// Get
	Get(id int) (*Report, error)
	// 查询
	All(param ListParam) (data []*Report, total int, err error)
	// 添加
	Add(uid string, content string) error
	// 更新
	Update(id int, content string) error
	// 赞
	Applaud(id int, uid string) error
	// 统计
	Stat(start, end time.Time) (stat *ReportStatResponse, err error)
	// AllStatus
	StatusRecords(status Status) (data []*ReportUser, err error)
	// AddStatus
	AddStatus(uid string, year, week int, status Status) error
	// RemoveStatus
	RemoveStatus(id int) error
	// StatusRecordsWithUser
	StatusRecordsWithUser(status Status, uid string) (data []*ReportStat, err error)
}

type Status int

const (
	WRNormal Status = iota
	WRVacation
	WRIgnore
)

type Report struct {
	Id      int             `json:"id"`
	Uid     string          `json:"uid"`
	Name    string          `json:"name,omitempty" db:"-"`   // staff name
	Avatar  string          `json:"avatar,omitempty" db:"-"` // staff name
	Year    int             `json:"year" db:"iso_year"`
	Week    int             `json:"week" db:"iso_week"`
	Status  Status          `json:"status"` // 0代表正常，1代表休假
	Content json.RawMessage `json:"content"`
	UpCount int             `json:"upCount" db:"up_count"`
	Created time.Time       `db:"created" json:"created"`
	Updated time.Time       `db:"updated" json:"updated,omitempty"`
}

type ReportWithProblem struct {
	Id   int `json:"id"`
	Year int `json:"year" db:"iso_year"`
	Week int `json:"week" db:"iso_week"`
	// Status    int64  `json:"status"` // 0代表正常，1代表休假
	Content   string    `json:"content"`
	ProblemId int       `json:"problem_id"`
	Problem   string    `json:"problem"`
	Created   time.Time `db:"created" json:"created"`
	Updated   time.Time `db:"updated,nullempty" json:"updated,omitempty"`
}

type ReportStat struct {
	Id      int       `json:"id,omitempty"`
	Uid     string    `json:"uid,omitempty"`
	Year    int       `json:"year" db:"iso_year"`
	Week    int       `json:"week" db:"iso_week"`
	Status  Status    `json:"status"` //0正常，1休假，2表示忽略该用户所有或特定周报
	Created time.Time `db:"created,omitempty" json:"created"`
}

type ReportUser struct {
	Id      int       `json:"id,omitempty"`
	Uid     string    `json:"uid"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// 周报统计参数
type ReportStatParam struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// 周报统计返回数据
type ReportStatResponse struct {
	Commited []*ReportStat `json:"commited"`
	All      []*ReportUser `json:"all"`
	Ignores  []*ReportUser `json:"ignores"`
}
