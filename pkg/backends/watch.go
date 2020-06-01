package backends

import (
	schema "github.com/liut/staffio-backend/model"
	"github.com/liut/staffio/pkg/models/team"
)

var _ team.WatchStore = (*watchStore)(nil)

type watchStore struct {
	ss schema.PeopleStore
}

func (s *watchStore) Gets(uid string) team.Butts {
	var data team.Butts

	err := withDbQuery(func(db dber) error {
		return db.Select(&data, "SELECT watching, created FROM staff_watch WHERE uid = $1", uid)
	})

	n := len(data)
	if err != nil {
		logger().Infow("watch gets fail", "err", err)
	} else if n > 0 {
		spec := &schema.Spec{UIDs: data.UIDs()}
		for _, staff := range s.ss.All(spec) {
			for i := 0; i < n; i++ {
				if staff.UID == data[i].UID {
					data[i].Name = staff.GetCommonName()
					data[i].Avatar = staff.AvatarURI()
				}
			}
		}
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
