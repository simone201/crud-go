package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	InitDB() *sql.DB
	CloseDB() error
	GetDB() *sql.DB
}

var Instance Database
