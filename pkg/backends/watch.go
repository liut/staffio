package backends

import (
	"github.com/liut/staffio/pkg/models/team"
)

var _ team.WatchStore = (*watchStore)(nil)

type watchStore struct {
}

func (s *watchStore) Gets(uid string) []string {
	var data []string

	err := withDbQuery(func(db dber) error {
		return db.Select(&data, "SELECT watching FROM staff_watch WHERE uid = $1", uid)
	})

	if err != nil {
		logger().Infow("watch gets fail", "err", err)
	}

	return data
}

func (s *watchStore) Watch(uid, watching string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		var id int
		err = db.Get(&id, "SELECT id FROM staff_watch WHERE uid = $1 AND watching = $2", uid, watching)
		if err == nil && id > 0 {
			logger().Infow("uid already watch", "uid", uid, "watch", watching)
			return nil
		}
		_, err = db.Exec("INSERT INTO staff_watch(uid, watching) VALUES ($1, $2)", uid, watching)
		return
	})
}

func (s *watchStore) Unwatch(uid, watching string) error {
	return withTxQuery(func(db dbTxer) (err error) {
		_, err = db.Exec("DELETE FROM staff_watch WHERE uid = $1 AND watching = $2", uid, watching)
		return
	})
}
