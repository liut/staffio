package backends

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/liut/staffio/pkg/models/weekly"
)

type teamStore struct{}

// Get
func (s *teamStore) Get(id int) (obj *weekly.Team, err error) {
	obj = &weekly.Team{}
	err = withDbQuery(func(db dber) error {
		return db.Get(obj,
			"SELECT id, name, leader, members, created FROM teams WHERE id = $1",
			id)
	})

	return
}

// 查询
func (s *teamStore) All(role weekly.TeamRoleType) (data []*weekly.Team, err error) {

	err = withDbQuery(func(db dber) error {
		data = make([]*weekly.Team, 0)
		switch role {
		case weekly.RoleMember:
			return db.Select(&data, `SELECT t.id, name, t.leader, members, tm.created, tm.uid as staff_uid
				FROM teams t JOIN team_member tm ON tm.team_id = t.id`)
		case weekly.RoleManager:
			return db.Select(&data, `SELECT t.id, name, t.leader, members, tm.created, tm.leader as staff_uid
				FROM teams t JOIN team_leader tm ON tm.team_id = t.id`)
		default:
			return db.Select(&data, "SELECT id, name, leader, members, created FROM teams")
		}

	})
	return
}

func (s *teamStore) Store(t *weekly.Team) error {
	return withTxQuery(func(db dbTxer) (err error) {
		if t.ID < 1 {
			var id int
			err = db.Get(&id, "INSERT INTO teams(name, leaders, members) VALUES($1, $2, $3) RETURNING id",
				t.Name, t.Leaders, t.Members)
			if err != nil {
				log.Printf("insert net team ERR %s", err)
			} else {
				log.Printf("insert new team id %d", id)
				t.ID = id
			}

		} else {
			var created time.Time
			err = db.Get(&created, "SELECT created FROM teams WHERE id = $1", t.ID)
			if err == nil {
				_, err = db.Exec(`UPDATE teams SET
					(name, leaders, members, updated) = ($1, $2, $3, CURRENT_TIMESTAMP) WHERE id = $4`,
					t.Name, t.Leaders, t.Members, t.ID)
				log.Printf("changed team %d %q %q", t.ID, t.Name, t.Leaders)
			} else if err == ErrNoRows {
				_, err = db.Exec("INSERT INTO teams(id, name, leaders, members) VALUES($1, $2, $3, $4)",
					t.ID, t.Name, t.Leaders, t.Members)
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
		_, err = db.Exec("DELETE FROM teams WHERE id = $1", id)
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
