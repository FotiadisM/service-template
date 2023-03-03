package db

import (
	"context"
	"database/sql"

	"github.com/FotiadisM/mock-microservice/pkg/logger"
	"go.uber.org/zap"
)

type TxFn func(tx *sql.Tx) (err error)

func WithTransaction(ctx context.Context, db *sql.DB, fn TxFn) error {
	return WithConfiguredTransaction(ctx, db, nil, fn)
}

func WithConfiguredTransaction(ctx context.Context, db *sql.DB, options *sql.TxOptions, fn TxFn) error {
	log := logger.FromContext(ctx)

	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		log.Error("failed to start transaction", zap.Error(err))
		return err
	}

	defer func() {
		p := recover()
		if p != nil {
			log.Info("recovered from panic, rolling back transaction and panicking again")

			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", zap.Error(err))
			}

			panic(p)
		}
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", zap.Error(err))
			}
		}
		if err = tx.Commit(); err != nil {
			log.Error("failed to commit transaction", zap.Error(err))
		}
	}()

	err = fn(tx)

	return err
}
