package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

type ActionData struct {
	Operation string `db:"operation"`
	Time      int64  `db:"time"`
	UID       int    `db:"uid"`
	Name      string `db:"name"`
	TID       int    `db:"tid"`
}

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "./data.db?_txlock=IMMEDIATE&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS actions (id PRIMARY KEY
		, operation TEXT NOT NULL
		, time INTEGER NOT NULL
		, uid INTEGER NOT NULL
		, name TEXT NOT NULL
		, tid INTEGER NOT NULL
		)`)
	if err != nil {
		panic(err)
	}
}

func Save(data ActionData) error {
	_, err := db.NamedExec(`INSERT INTO actions (operation, time, uid, name, tid) VALUES (:operation, :time, :uid, :name, :tid)`, data)
	if err != nil {
		return fmt.Errorf("Save: %w", err)
	}
	return nil
}
