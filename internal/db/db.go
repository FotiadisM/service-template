package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	// PostgreSQL databse driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/db/repository"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type TxFn func(db DB) (err error)

type DB interface {
	repository.Querier
	Ping(ctx context.Context) error
	WithTx(ctx context.Context, fn TxFn) error
	WithConfiguredTx(ctx context.Context, options *sql.TxOptions, fn TxFn) error
}

type db struct {
	*sql.DB
	*repository.Queries
}

func New(config config.DB) (DB, error) {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Host, config.Port, config.Username, config.Password, config.Database)
	for k, v := range config.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	conn, err := otelsql.Open("pgx", str)
	if err != nil {
		return nil, err
	}

	if config.MaxOpenConns > 0 {
		conn.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		conn.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		conn.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	db := &db{
		DB:      conn,
		Queries: repository.New(conn),
	}

	return db, nil
}

func (s *db) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *db) WithTx(ctx context.Context, fn TxFn) error {
	return s.WithConfiguredTx(ctx, nil, fn)
}

func (s *db) WithConfiguredTx(ctx context.Context, options *sql.TxOptions, fn TxFn) error {
	log := ilog.FromContext(ctx)

	tx, err := s.DB.BeginTx(ctx, options)
	if err != nil {
		log.Error("failed to begin transaction", ilog.Err(err.Error()))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error("recovered from panic, rolling back transaction and panicking again")
			if err = tx.Rollback(); err != nil {
				log.Error("failed to roll back transaction", ilog.Err(err.Error()))
			}
			panic(p)
		}
	}()

	err = fn(&db{
		DB:      s.DB,
		Queries: s.Queries.WithTx(tx),
	})
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Error("failed to roll back transaction", ilog.Err(err.Error()))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
