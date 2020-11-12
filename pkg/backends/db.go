package backends

import (
	"database/sql"
	"errors"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // ok

	"github.com/jmoiron/sqlx"
)

var (
	ErrEmptyVal = errors.New("empty value")
	dbError     = errors.New("database error")
	ErrNotFound = errors.New("Not Found")
	valueError  = errors.New("value error")
	dbc         *sqlx.DB
	dbDSN       string

	once         sync.Once
	quitC        chan struct{}
	pingInterval = 90 * time.Second

	ErrNoRows = sql.ErrNoRows
)

type dber interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type dbTxer interface {
	dber
	Rollback() error
	Commit() error
}

func init() {
	dbDSN = envOr("STAFFIO_BACKEND_DSN", "postgres://staffio@localhost/staffio?sslmode=disable")

	quitC = make(chan struct{})
}

func SetDSN(dsn string) {
	if dsn != "" {
		logger().Debugw("set db dsn", "dsn", dsn)
		dbDSN = dsn
		openDb()
	}
}

func openDb() *sqlx.DB {
	db, err := sqlx.Open("postgres", dbDSN)
	if err != nil {
		logger().Errorw("open db fail", "err", err)
	}

	return db
}

func closeDb() {
	if dbc != nil {
		err := dbc.Close()
		if err != nil {
			logger().Warnw("close db fail", "err", err)
		}
	}
}

func getDb() *sqlx.DB {
	if dbc == nil {
		once.Do(func() {
			dbc = openDb()
			go reap(pingInterval, pingDb, quitC)
		})
	}

	return dbc
}

func pingDb() error {
	err := dbc.Ping()
	if err != nil {
		logger().Infow("ping db fail", "err", err)
		dbc = openDb()
	}
	return err
}

// reap with special action at set intervals.
func reap(interval time.Duration, cf func() error, quit <-chan struct{}) {
	logger().Debugw("starting reaper", "interval", interval)
	ticker := time.NewTicker(interval)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-quit:
			// Handle the quit signal.
			return
		case <-ticker.C:
			// Execute function of clean.
			if err := cf(); err != nil {
				logger().Infow("reap fail", "err", err)
			}
		}
	}
}

func withDbQuery(query func(db dber) error) error {
	db := getDb()
	// defer db.Close()
	if err := query(db); err != nil {
		if err == sql.ErrNoRows {
			logger().Infow("db query fail", "err", err)
			return ErrNotFound
		}
		logger().Warnw("db query fail", "err", err)
		return dbError
	}
	return nil
}

func withTxQuery(query func(tx dbTxer) error) error {

	db := getDb()
	// defer db.Close()

	tx, err := db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := query(tx); err != nil {
		tx.Rollback()
		logger().Warnw("tx query fail", "err", err)
		return dbError
	}
	tx.Commit()
	return nil
}

func inArray(k string, fields []string) bool {
	for _, sf := range fields {
		if k == sf {
			return true
		}
	}
	return false
}

func envOr(key, dft string) string {
	v := os.Getenv(key)
	if v == "" {
		return dft
	}
	return v
}
