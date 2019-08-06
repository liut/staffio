package backends

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	ErrEmptyVal = errors.New("empty value")
	dbError     = errors.New("database error")
	ErrNotFound = errors.New("Not Found")
	valueError  = errors.New("value error")
	dbc         *sqlx.DB
	dbDSN       string

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
	if s, exists := os.LookupEnv("STAFFIO_BACKEND_DSN"); exists && s != "" {
		dbDSN = s
	} else {
		dbDSN = "postgres://staffio@localhost/staffio?sslmode=disable"
	}
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
		dbc = openDb()
		return dbc
	}

	if err := dbc.Ping(); err != nil {
		dbc = openDb()
	}

	return dbc
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
