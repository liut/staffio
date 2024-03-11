package backends

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/liut/staffio/pkg/models/weekly"
)

var _ weekly.Store = (*weeklyStore)(nil)

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
func (s *weeklyStore) All(spec weekly.ReportsSpec) (data weekly.Reports, total int, err error) {
	logger().Debugw("weeklyStore.all", "spec", spec)
	var where string
	bind := []any{}
	if len(spec.UIDs) > 0 {
		var args []any
		where, args, err = sqlx.In("WHERE r.uid IN (?)", spec.UIDs)
		if err != nil {
			return
		}
		bind = append(bind, args...)
	} else if spec.UID != "" {
		where = " WHERE r.uid = ?"
		bind = append(bind, spec.UID)
	}

	if spec.TeamID > 0 {
		if where == "" {
			where = " WHERE team_id = ?"
		} else {
			where += " AND team_id = ?"
		}
		bind = append(bind, spec.TeamID)
	}

	if err = withDbQuery(func(db dber) error {
		where = sqlx.Rebind(sqlx.DOLLAR, where)
		logger().Debugw("query weekly_report", "where", where, "bind", bind)
		return db.Get(&total, "SELECT COUNT(DISTINCT r.id) FROM weekly_report r "+
			"LEFT JOIN team_member tm ON tm.uid = r.uid "+
			where, bind...)
	}); err != nil {
		logger().Infow("query weekly report fail", "where", where, "err", err)
		return
	}

	str := `SELECT DISTINCT r.id, r.uid, iso_year, iso_week, content, r.created, r.updated, r.up_count
	   FROM weekly_report r LEFT JOIN team_member tm ON tm.uid = r.uid ` +
		where +
		spec.Sort.Sql() + spec.Pager.Sql()

	data = make(weekly.Reports, 0)
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
			_, err = db.Exec("UPDATE weekly_report SET content = $1, updated = now() WHERE id = $2", content, id)
			return err
		}
		return err
	})
}

// 更新
func (s *weeklyStore) Update(id int, content string) (err error) {
	qs := func(db dbTxer) error {
		if id > 0 {
			_, err = db.Exec("UPDATE weekly_report SET content = $1, updated = now() WHERE id = $2", content, id)
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
	year1, _ := start.ISOWeek()
	year2, week2 := end.ISOWeek()
	qs := func(db dber) error {
		return db.Select(&rsr.Commited,
			// TODO: 使用 created 来作为筛选条件， 如果将来加上了补交功能，则这块儿需要更改
			// 如果 report 和 status 两张表有 uid、year、week 三个字段都想同的数据，以 status 为准
			"SELECT id, uid, iso_year,iso_week, 0 AS status, created FROM weekly_report "+
				"WHERE created >=$1 AND created <=$2 "+
				"AND NOT EXISTS(SELECT uid, iso_year, iso_week FROM weekly_status "+
				"WHERE weekly_report.uid=weekly_status.uid "+
				"AND weekly_report.iso_year=weekly_status.iso_year "+
				"AND weekly_report.iso_week=weekly_status.iso_week) UNION "+
				"SELECT id, uid,iso_year,iso_week,status, created FROM weekly_status "+
				"WHERE iso_year>0 AND iso_week>0 AND (iso_year >= $5 AND iso_year<=$3 OR (iso_year=$3 AND iso_week<=$4))",
			start,
			end,
			year2,
			week2,
			year1,
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
func (s *weeklyStore) AddStatus(uid string, status weekly.Status, year int, weeks ...int) error {
	qs := func(db dbTxer) (err error) {
		for _, week := range weeks {
			var id int
			err = db.Get(&id, "SELECT id FROM weekly_status WHERE uid = $1 AND iso_year = $2 AND iso_week = $3",
				uid, year, week)
			if err == nil {
				_, err = db.Exec("UPDATE weekly_status SET status = $1 WHERE id = $2", status, id)
			} else if err == ErrNoRows {
				_, err = db.Exec("INSERT INTO weekly_status(uid, iso_year, iso_week, status) VALUES($1, $2, $3, $4)",
					uid, year, week, status)
			}
			if err != nil {
				return
			}
		}

		return
	}
	return withTxQuery(qs)
}

// RemoveStatus
func (s *weeklyStore) RemoveStatus(ids ...int) error {
	return withTxQuery(func(db dbTxer) (err error) {
		for _, id := range ids {
			_, err = db.Exec("DELETE FROM weekly_status WHERE id = $1 ",
				id)
			if err != nil {
				return
			}
		}

		return
	})
}
