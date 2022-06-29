package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type RollbackFunc func(tx *sql.Tx)

func initBackgroundTransaction(db Database, ctx context.Context, timeout time.Duration) (*sql.Tx, context.Context, context.CancelFunc, RollbackFunc, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout*time.Second)

	tx, err := db.GetDB().BeginTx(ctx, nil)
	if err != nil {
		cancelFunc()
		return nil, nil, nil, nil, err
	}

	rollbackFunc := func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Print(err)
		}
	}

	return tx, ctx, cancelFunc, rollbackFunc, nil
}
