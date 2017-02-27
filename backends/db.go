package backends

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	. "lcgc/platform/staffio/settings"
)

var (
	dbError    = errors.New("database error")
	valueError = errors.New("value error")
	dbc        *sqlx.DB
)

const (
	ASCENDING  = 1
	DESCENDING = -1
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

func openDb() *sqlx.DB {
	log.Printf("PGHOST: %s", os.Getenv("PGHOST"))
	db, err := sqlx.Open("postgres", Settings.Backend.DSN)
	if err != nil {
		log.Fatalf("open db error: %s", err)
	}

	return db
}

func closeDb() {
	if dbc != nil {
		err := dbc.Close()
		if err != nil {
			log.Printf("closing db error: %s", err)
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
		log.Printf("db query ERR: %s", err)
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
		log.Printf("tx query ERR: %s", err)
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
