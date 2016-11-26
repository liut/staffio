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
	return withDbQuery(func(db *sql.DB) (err error) {
		_, err = db.Exec("DELETE FROM oauth_authorization_code WHERE created < $1", now.Add(-time.Second*AuthorizationExpiration))
		if err != nil {
			log.Printf("clean authorize ERR %s", err)
			return
		}
		_, err = db.Exec("DELETE FROM oauth_access_token WHERE created < $1", now.Add(-time.Second*AccessExpiration))
		if err != nil {
			log.Printf("clean access ERR %s", err)
			return
		}

		_, err = db.Exec("DELETE FROM http_sessions WHERE expires_on < now()")
		if err != nil {
			log.Printf("clean sessions ERR %s", err)
			return
		}
		// log.Printf("Cleanup done at %s", now)
		return
	})
}

// a+x>b
