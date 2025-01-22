package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	// PostgreSQL databse driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/FotiadisM/mock-microservice/internal/services/auth/v1/queries"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type DB struct {
	DB *sql.DB
	*queries.Queries
}

func New(config config.DB) (*DB, error) {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Host, config.Port, config.Username, config.Password, config.Database)
	for k, v := range config.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	conn, err := otelsql.Open("pgx", str)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
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

	db := &DB{
		DB:      conn,
		Queries: queries.New(conn),
	}

	return db, nil
}

func NewFromDBConn(conn *sql.DB) (*DB, error) {
	db := &DB{
		DB:      conn,
		Queries: queries.New(conn),
	}

	return db, nil
}

type TxFn func(db *DB) (err error)

func WithTx(ctx context.Context, db *DB, fn TxFn) error {
	return WithConfiguredTx(ctx, db, nil, fn)
}

func WithConfiguredTx(ctx context.Context, db *DB, options *sql.TxOptions, fn TxFn) error {
	log := ilog.FromContext(ctx)

	tx, err := db.DB.BeginTx(ctx, options)
	if err != nil {
		log.Error("failed to begin transaction", ilog.Err(err))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error("recovered from panic, rolling back transaction and panicking again")
			if err = tx.Rollback(); err != nil {
				log.Error("failed to roll back transaction", ilog.Err(err))
			}
			panic(p)
		}
	}()

	err = fn(&DB{
		DB:      db.DB,
		Queries: db.Queries.WithTx(tx),
	})
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Error("failed to roll back transaction", ilog.Err(err))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
