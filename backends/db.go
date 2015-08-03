package backends

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	. "tuluu.com/liut/staffio/settings"
)

var (
	dbError    = errors.New("database error")
	valueError = errors.New("value error")
	dbc        *sql.DB
)

const (
	ASCENDING  = 1
	DESCENDING = -1
)

func openDb() *sql.DB {
	db, err := sql.Open("postgres", Settings.Backend.DSN)
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

func getDb() *sql.DB {
	if dbc == nil {
		dbc = openDb()
		return dbc
	}

	if err := dbc.Ping(); err != nil {
		dbc = openDb()
	}

	return dbc
}

func withDbQuery(query func(db *sql.DB) error) error {
	db := getDb()
	// defer db.Close()
	if err := query(db); err != nil {
		log.Printf("db query error: %s", err)
		return dbError
	}
	return nil
}

func withTxQuery(query func(tx *sql.Tx) error) error {

	db := getDb()
	// defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := query(tx); err != nil {
		tx.Rollback()
		log.Printf("tx query error: %s", err)
		return dbError
	}
	tx.Commit()
	return nil
}

func inSortable(k string, fields []string) bool {
	for _, sf := range fields {
		if k == sf {
			return true
		}
	}
	return false
}
