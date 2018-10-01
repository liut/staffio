package backends

import (
	"fmt"
	"time"

	"github.com/liut/staffio/pkg/models/weekly"
)

type weeklyStore struct {
}

// Get
func (s *weeklyStore) Get(id int) (obj *weekly.Report, err error) {
	obj = &weekly.Report{}
	qs := func(db dber) error {
		return db.Get(obj,
			"SELECT id, uid, iso_year, iso_week, content, created FROM weekly_report WHERE id = $1",
			id)
	}
	err = withDbQuery(qs)

	return
}

// 查询
func (s *weeklyStore) All(param weekly.ListParam) (data []*weekly.Report, total int, err error) {

	var where string
	bind := []interface{}{}
	if param.GroupId > 0 {
		where = " WHERE team_id = $1"
		bind = append(bind, param.GroupId)
	}
	if param.Uid != "" {
		if where == "" {
			where = " WHERE r.uid = $1"
		} else {
			where += " AND r.uid = $2"
		}
		bind = append(bind, param.Uid)
	}

	if err = withDbQuery(func(db dber) error {
		return db.Get(&total, "SELECT COUNT(r.id) FROM weekly_report r "+
			"LEFT JOIN team_member tm ON tm.uid = r.uid "+
			where, bind...)
	}); err != nil {
		return
	}

	str := `SELECT r.id, r.uid, iso_year, iso_week, content, r.created, r.updated, r.up_count
	   FROM weekly_report r LEFT JOIN team_member tm ON tm.uid = r.uid ` +
		where +
		param.Sort.Sql() + param.Pager.Sql()

	data = make([]*weekly.Report, 0)
	qs := func(db dber) error {
		return db.Select(&data, str, bind...)
	}

	if err = withDbQuery(qs); err != nil {
		return
	}

	return
}

// 添加
func (s *weeklyStore) Add(uid string, content string) error {
	now := time.Now()
	year, week := now.ISOWeek()

	return withTxQuery(func(db dbTxer) error {
		var id int
		err := db.Get(&id,
			"SELECT id FROM weekly_report WHERE uid = $1 AND iso_year = $2 AND iso_week = $3", uid, year, week)
		if err == ErrNoRows {
			return db.Get(&id,
				"INSERT INTO weekly_report (uid, iso_year, iso_week, content) VALUES ($1,$2,$3,$4)"+
					" RETURNING id", uid, year, week, content)
		}
		if id > 0 {
			_, err = db.Exec("UPDATE weekly_report SET content = $1 WHERE id = $2", content, id)
			return err
		}
		return err
	})
}

// 更新
func (s *weeklyStore) Update(id int, content string) (err error) {
	qs := func(db dbTxer) error {
		if id > 0 {
			_, err = db.Exec("UPDATE weekly_report SET content = $1 WHERE id = $2", content, id)
			return err
		}
		return fmt.Errorf("invalid id value %d", id)
	}
	err = withTxQuery(qs)
	return
}

// 赞
func (s *weeklyStore) Applaud(id int, uid string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		_, err = db.Exec("INSERT INTO weekly_report_up (report_id, uid) VALUES ($1,$2)",
			id, uid)
		if err == nil {
			_, err = db.Exec(
				`UPDATE weekly_report wr SET up_count = (
				SELECT COUNT(id) FROM weekly_report_up wra WHERE wra.report_id = wr.id)`)
		}
		return
	})
}

// 统计
func (s *weeklyStore) Stat(start, end time.Time) (rsr *weekly.ReportStatResponse, err error) {
	rsr = &weekly.ReportStatResponse{}
	rsr.Commited = []*weekly.ReportStat{}
	year, week := end.ISOWeek()
	qs := func(db dber) error {
		return db.Select(&rsr.Commited,
			// TODO: 使用create_at来作为筛选条件，如果将来加上了补交功能，则这块儿需要更改
			// 如果report和status两张表有user_id、year、week三个字段都想同的数据，以status为准
			"SELECT id, uid, iso_year,iso_week, 0 AS status, created FROM weekly_report "+
				"WHERE created >=$1 AND created <=$2 "+
				"AND NOT EXISTS(SELECT uid, iso_year, iso_week FROM weekly_status "+
				"WHERE weekly_report.uid=weekly_status.uid "+
				"AND weekly_report.iso_year=weekly_status.iso_year "+
				"AND weekly_report.iso_week=weekly_status.iso_week) UNION "+
				"SELECT id, uid,iso_year,iso_week,status, created FROM weekly_status "+
				"WHERE iso_year>0 AND iso_week>0 AND (iso_year<$3 OR (iso_year=$3 AND iso_week<=$4))",
			start,
			end,
			year,
			week,
		)
	}
	err = withDbQuery(qs)
	rsr.All = []*weekly.ReportUser{}
	// rsr.Ignores, err = s.StatusRecords(weekly.WRIgnore)
	return
}

// StatusRecords
func (s *weeklyStore) StatusRecords(status weekly.Status) (data []*weekly.ReportUser, err error) {
	data = []*weekly.ReportUser{}
	qs := func(db dber) error {
		return db.Select(&data,
			"SELECT id, uid, created FROM weekly_status "+
				"WHERE status=$1", status)

	}
	err = withDbQuery(qs)
	return
}

// StatusRecordsWithUser
func (s *weeklyStore) StatusRecordsWithUser(status weekly.Status, uid string) (data []*weekly.ReportStat, err error) {
	data = []*weekly.ReportStat{}
	qs := func(db dber) error {
		return db.Select(&data,
			"SELECT id, iso_year, iso_week, status from weekly_status "+
				"WHERE uid = $1 AND status=$2 AND iso_week > 0 ORDER BY iso_year, iso_week", uid, status)
	}
	err = withDbQuery(qs)
	return
}

// AddStatus
func (s *weeklyStore) AddStatus(uid string, year, week int, status weekly.Status) error {
	qs := func(db dbTxer) error {
		var id int
		err := db.Get(&id, "SELECT id FROM weekly_status WHERE uid = $1 AND iso_year = $2 AND iso_week = $3",
			uid, year, week)
		if err == nil {
			_, err = db.Exec("UPDATE weekly_status SET status = $1 WHERE id = $2", status, id)
			return err
		}
		_, err = db.Exec("INSERT INTO weekly_status(uid, iso_year, iso_week, status) VALUES($1, $2, $3, $4)",
			uid, year, week, status)
		return err
	}
	return withTxQuery(qs)
}

// RemoveStatus
func (s *weeklyStore) RemoveStatus(id int) error {
	return withTxQuery(func(db dbTxer) error {
		_, err := db.Exec("DELETE FROM weekly_status WHERE id = $1 ",
			id)
		return err
	})
}
