package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"
)

const dbFileName = "data.db"

type SqliteDatabase struct {
	DB *sql.DB
}

func (sd *SqliteDatabase) InitDB() *sql.DB {
	if _, err := os.Stat(dbFileName); errors.Is(err, os.ErrNotExist) {
		_, fileErr := os.Create(dbFileName)
		if fileErr != nil {
			log.Fatal(fileErr)
		}
	}

	dbInstance, dbErr := sql.Open("sqlite3", dbFileName)
	if dbErr != nil {
		log.Fatal(dbErr)
	} else {
		sd.DB = dbInstance
	}

	createTables(dbInstance)

	return dbInstance
}

func (sd *SqliteDatabase) CloseDB() error {
	return sd.DB.Close()
}

func (sd *SqliteDatabase) GetDB() *sql.DB {
	return sd.DB
}

func createTables(db *sql.DB) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	createPeopleTable(db, ctx)
}

func createPeopleTable(db *sql.DB, ctx context.Context) {
	createPeopleQuery := "CREATE TABLE IF NOT EXISTS `people` (\n" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, \n" +
		"`name` VARCHAR(64) NULL, \n" +
		"`birth` DATETIME NULL, \n" +
		"`createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, \n" +
		"`updatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP \n" +
		")"

	_, err := createInContext(db, ctx, createPeopleQuery)
	if err != nil {
		log.Fatal(err)
	}
}
