package models

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database interface {
	PrepareStmt(string) (*sql.Stmt, error)
}

type MysqlDB struct {
	db *sql.DB
}

func NewMysqlDB(conn string) (*MysqlDB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MysqlDB{db: db}, nil
}

func (m *MysqlDB) PrepareStmt(query string) (*sql.Stmt, error) {
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}
