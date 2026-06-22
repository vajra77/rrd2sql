package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Sql3Db struct {
	filename string
	db       *sql.DB
}

func NewSql3Db(filename string) (*Sql3Db, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open DB: %v", err)
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS traffic (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        tstamp 	INTEGER,
        avg_in 	FLOAT,
        avg_out FLOAT,
        max_in 	FLOAT,
        max_out FLOAT
    );
    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("unable to create table: %v", err)
	}

	return &Sql3Db{filename, db}, nil
}

func (s *Sql3Db) Close() error {
	err := s.db.Close()
	if err != nil {
		fmt.Printf("ERROR: unable to close DB: %v", err)
	}
}

func (s *Sql3Db) WriteResult(result Result) error {
	stmt, err := s.db.Prepare("INSERT INTO traffic (tstamp, avg_in, avg_out, max_in, max_out) VALUES (?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(result.Timestamp, result.AvgIn, result.AvgOut, result.MaxIn, result.MaxOut)
	if err != nil {
		return fmt.Errorf("unable to write to DB: %v", err)
	}
	return nil
}
