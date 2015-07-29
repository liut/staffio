package backends

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	. "tuluu.com/liut/staffio/settings"
)

var (
	dbc *sql.DB
)

func openDb() *sql.DB {
	db, err := sql.Open("postgres", Settings.Backend.DSN)
	if err != nil {
		log.Fatalf("open db error: %s", err)
	}
	return db
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

func withTxQuery(query func(tx *sql.Tx) error) error {

	db := getDb()
	defer db.Close()

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
