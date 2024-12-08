package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// PostgreSQL databse driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/store/repository"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type TxFn func(store Store) (err error)

type Store interface {
	repository.Querier
	Ping(ctx context.Context) error
	WithTx(ctx context.Context, fn TxFn) error
	WithConfiguredTx(ctx context.Context, options *sql.TxOptions, fn TxFn) error
}

type store struct {
	*sql.DB
	*repository.Queries
}

func New(config config.DB) (Store, error) {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Host, config.Port, config.Username, config.Password, config.Database)
	for k, v := range config.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	db, err := sql.Open("pgx", str)
	if err != nil {
		return nil, err
	}

	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	store := &store{
		DB:      db,
		Queries: repository.New(db),
	}

	return store, nil
}

func (s *store) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *store) WithTx(ctx context.Context, fn TxFn) error {
	return s.WithConfiguredTx(ctx, nil, fn)
}

func (s *store) WithConfiguredTx(ctx context.Context, options *sql.TxOptions, fn TxFn) error {
	log := ilog.FromContext(ctx)

	tx, err := s.DB.BeginTx(ctx, options)
	if err != nil {
		log.Error("failed to begin transaction", ilog.Err(err.Error()))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error("recovered from panic, rolling back transaction and panicking again")
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", ilog.Err(err.Error()))
			}
			panic(p)
		}
	}()

	err = fn(&store{
		DB:      s.DB,
		Queries: s.Queries.WithTx(tx),
	})
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to roll back transaction: %w", txErr))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
