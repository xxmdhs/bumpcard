package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ActionData struct {
	Operation string `db:"operation" json:"operation"`
	Time      int64  `db:"time" json:"time"`
	UID       int    `db:"uid" json:"uid"`
	Name      string `db:"name" json:"name"`
	TID       int    `db:"tid" json:"tid"`
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

func (db *DB) Del(tid int) error {
	_, err := db.db.Exec(`DELETE FROM actions WHERE tid = ?`, tid)
	if err != nil {
		return fmt.Errorf("Del: %w", err)
	}
	return nil
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

func (db *DB) GetForUID(uid int) ([]ActionData, error) {
	var data []ActionData
	err := db.db.Select(&data, `SELECT operation, time, uid, name, tid FROM actions WHERE uid = ?`, uid)
	if err != nil {
		return nil, fmt.Errorf("GetForUID: %w", err)
	}
	return data, nil
}
