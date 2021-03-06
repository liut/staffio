package backends

import (
	"fmt"
	"log"
	"time"
)

const (
	authorizationExpiration = 60 * 15
	accessExpiration        = 60 * 60 * 24
	passwordExpiration      = 60 * 120
	sessionExpiration       = 60 * 30
)

// Cleanup 清理过期的数据
func Cleanup() (err error) {
	now := time.Now()

	err = deleteWithEnd("oauth_authorization_code", "created", now.Add(-time.Second*authorizationExpiration))
	if err != nil {
		return
	}
	err = deleteWithEnd("oauth_access_token", "created", now.Add(-time.Second*accessExpiration))
	if err != nil {
		return
	}
	err = deleteWithEnd("password_reset", "created", now.Add(-time.Second*passwordExpiration))
	if err != nil {
		return
	}
	err = deleteWithEnd("http_sessions", "expires_on", now.Add(-time.Second*sessionExpiration))
	if err != nil {
		return
	}
	return
}

func deleteWithEnd(name, field string, end time.Time) error {
	return withDbQuery(func(db dber) error {
		qs := fmt.Sprintf("DELETE FROM %s WHERE %s < $1", name, field)
		res, err := db.Exec(qs, end)
		if err != nil {
			log.Printf("clean %q ERR %s", name, err)
			return err
		}
		count, _ := res.RowsAffected()
		log.Printf("clean %q: %d affected", name, count)
		return nil
	})
}

// a+x>b
