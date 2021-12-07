package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ActionData struct {
	Operation string `db:"operation"`
	Time      int64  `db:"time"`
	UID       int    `db:"uid"`
	Name      string `db:"name"`
	TID       int    `db:"tid"`
}

type DB struct {
	db *sqlx.DB
}

func NewSql(filename string) (*DB, error) {
	var err error
	db := &DB{}
	db.db, err = sqlx.Connect("sqlite3", "./data.db?_txlock=IMMEDIATE&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("newSql: %w", err)
	}
	_, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS actions (id INTEGER PRIMARY KEY AUTOINCREMENT
		, operation TEXT NOT NULL
		, time INTEGER NOT NULL
		, uid INTEGER NOT NULL
		, name TEXT NOT NULL
		, tid INTEGER NOT NULL
		)`)
	if err != nil {
		return nil, fmt.Errorf("newSql: %w", err)
	}
	return db, nil
}

func (db *DB) Save(data ActionData) error {
	_, err := db.db.NamedExec(`INSERT INTO actions (operation, time, uid, name, tid) VALUES (:operation, :time, :uid, :name, :tid)`, data)
	if err != nil {
		return fmt.Errorf("Save: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
