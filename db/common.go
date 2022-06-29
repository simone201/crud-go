package db

import "database/sql"

type Database interface {
	InitDB() *sql.DB
	CloseDB() error
	GetDB() *sql.DB
}

var Instance Database
