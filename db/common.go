package db

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Database interface {
	InitDB() *sql.DB
	CloseDB() error
	GetDB() *sql.DB
}

var Instance Database

func createInContext(db *sql.DB, ctx context.Context, query string) (sql.Result, error) {
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating table", err)
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return nil, err
	}
	log.Printf("Rows affected when creating table: %d", rows)

	return res, nil
}
