package backends

import (
	"database/sql"
	"log"
	"time"
)

const (
	AuthorizationExpiration = 900
	AccessExpiration        = 86400
)

// 清理过期的数据
func Cleanup() error {
	now := time.Now()
	return withDbQuery(func(db dber) (err error) {
		var (
			r1, r2, r3 sql.Result
			c1, c2, c3 int64
		)
		r1, err = db.Exec("DELETE FROM oauth_authorization_code WHERE created < $1", now.Add(-time.Second*AuthorizationExpiration))
		if err != nil {
			log.Printf("clean authorize ERR %s", err)
			return
		}
		c1, _ = r1.RowsAffected()
		r2, err = db.Exec("DELETE FROM oauth_access_token WHERE created < $1", now.Add(-time.Second*AccessExpiration))
		if err != nil {
			log.Printf("clean access ERR %s", err)
			return
		}
		c2, _ = r2.RowsAffected()
		r3, err = db.Exec("DELETE FROM http_sessions WHERE expires_on < now()")
		if err != nil {
			log.Printf("clean sessions ERR %s", err)
			return
		}
		c3, _ = r3.RowsAffected()
		debug("Cleanup done at %s: %d, %d, %d", now, c1, c2, c3)
		return
	})
}

// a+x>b
