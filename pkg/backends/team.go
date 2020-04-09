package backends

import (
	"fmt"
	"strings"
	"time"

	"github.com/liut/staffio/pkg/models/team"
)

type teamStore struct{}

// Get
func (s *teamStore) Get(id int) (obj *team.Team, err error) {
	obj = &team.Team{}
	err = withDbQuery(func(db dber) error {
		return db.Get(obj,
			"SELECT id, name, leaders, members, created FROM teams WHERE id = $1",
			id)
	})

	return
}

func (s *teamStore) GetWithMember(uid string) (obj *team.Team, err error) {
	obj = new(team.Team)
	err = withDbQuery(func(db dber) error {
		var teamID int
		err := db.Get(&teamID, "SELECT team_id FROM team_member WHERE uid = $1", uid)
		if err == nil {
			return db.Get(obj, "SELECT * FROM teams WHERE id = $1", teamID)
		}
		return err
	})
	return
}

// 查询
func (s *teamStore) All(role team.RoleType) (data team.Teams, err error) {

	err = withDbQuery(func(db dber) error {
		data = make(team.Teams, 0)
		switch role {
		case team.RoleMember:
			return db.Select(&data, `SELECT t.id, name, leaders, members, tm.created, tm.uid as staff_uid
				FROM teams t JOIN team_member tm ON tm.team_id = t.id`)
		case team.RoleManager:
			return db.Select(&data, `SELECT t.id, name, leaders, members, tm.created, tm.leader as staff_uid
				FROM teams t JOIN team_leader tm ON tm.team_id = t.id`)
		default:
			return db.Select(&data, "SELECT id, name, leaders, members, created FROM teams ORDER BY id")
		}

	})
	return
}

func (s *teamStore) Store(t *team.Team) error {
	if t.Name == "" {
		return ErrEmptyVal
	}
	return withTxQuery(func(db dbTxer) (err error) {
		if t.ID < 1 {
			var id int
			if err = db.Get(&id, "SELECT id FROM teams WHERE name = $1", t.Name); err == nil {
				return
			}
			err = db.Get(&id, "INSERT INTO teams(name, parent_id, leaders, members) VALUES($1, $2, $3, $4) RETURNING id",
				t.Name, t.ParentID, t.Leaders, t.Members)
			if err != nil {
				logger().Infow("insert new team fail", "err", err)
			} else {
				logger().Infow("insert new team ok", "name", t.Name, "id", id)
				t.ID = id
			}

		} else {
			var created time.Time
			err = db.Get(&created, "SELECT created FROM teams WHERE id = $1", t.ID)
			if err == nil {
				_, err = db.Exec(`UPDATE teams SET
					(name, parent_id, leaders, members, updated) = ($1, $2, $3, $4, CURRENT_TIMESTAMP) WHERE id = $5`,
					t.Name, t.ParentID, t.Leaders, t.Members, t.ID)
			} else if err == ErrNoRows {
				db.Exec("DELETE FROM teams WHERE parent_id = $1 AND name = $2", t.ParentID, t.Name)
				_, err = db.Exec("INSERT INTO teams(id, name, parent_id, leaders, members) VALUES($1, $2, $3, $4, $5)",
					t.ID, t.Name, t.ParentID, t.Leaders, t.Members)
			}
			if err != nil {
				logger().Infow("store team fail", "team", t, "err", err)
			}
		}
		if err == nil {
			err = dbTeamAddMember(db, t.ID, t.Members)
		}
		return
	})
}

func (s *teamStore) Delete(id int) error {
	return withTxQuery(func(db dbTxer) (err error) {
		_, err = db.Exec("DELETE FROM team_leader WHERE team_id = $1", id)
		if err == nil {
			_, err = db.Exec("DELETE FROM team_member WHERE team_id = $1", id)
			if err == nil {
				_, err = db.Exec("DELETE FROM teams WHERE id = $1", id)
			}
		}

		return
	})
}

// Add members
func (s *teamStore) AddMember(id int, uids ...string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		return dbTeamAddMember(db, id, uids)
	})
}

func dbTeamAddMember(db dbTxer, id int, uids []string) (err error) {
	for _, uid := range uids {
		uid = strings.ToLower(uid)
		var existID int
		if db.Get(&existID, "SELECT id FROM team_member WHERE team_id = $1 AND uid = $2", id, uid) == nil {
			continue
		}
		_, err = db.Exec("INSERT INTO team_member(team_id, uid) VALUES($1, $2)", id, uid)
		if err != nil {
			return
		}
	}
	return
}

// Remove members
func (s *teamStore) RemoveMember(id int, uids ...string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		var arr []string
		var bind = []interface{}{id}
		for i, s := range uids {
			arr = append(arr, fmt.Sprintf("$%d", i+2))
			bind = append(bind, s)
		}
		_, err = db.Exec("DELETE FROM team_member WHERE team_id = $1 AND uid IN ("+strings.Join(arr, ",")+") ",
			bind...)
		if err != nil {
			logger().Infow("delete team member fail", "id", id, "uids", uids, "err", err)
		} else {
			logger().Infow("delete team member done", "id", id, "uids", uids)
		}
		return
	})
}

// Add Manager
func (s *teamStore) AddManager(id int, uid string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		uid = strings.ToLower(uid)
		var existID int
		if db.Get(&existID, "SELECT id FROM team_leader WHERE team_id = $1 AND leader = $2", id, uid) == nil {
			return
		}
		_, err = db.Exec("INSERT INTO team_leader(team_id, leader) VALUES($1, $2)", id, uid)
		return
	})
}

// Remove Manager
func (s *teamStore) RemoveManager(id int, uid string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		_, err = db.Exec("DELETE FROM team_leader WHERE team_id = $1 AND leader = $2",
			id, strings.ToLower(uid))
		return
	})
}
