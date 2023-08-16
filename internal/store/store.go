package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// PostgreSQL databse driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/FotiadisM/mock-microservice/internal/store/queries"
	"github.com/FotiadisM/mock-microservice/pkg/ilog"
)

type TxFn func(store Store) (err error)

type Store interface {
	queries.Querier
	Ping(ctx context.Context) error
	WithTx(ctx context.Context, fn TxFn) error
	WithConfiguredTx(ctx context.Context, options *sql.TxOptions, fn TxFn) error
}

type Config struct {
	Host            string            `env:"PSQL_HOST,default=localhost"`
	Port            int               `env:"PSQL_PORT,default=5432"`
	Username        string            `env:"PSQL_USER,default=local_user"`
	Password        string            `env:"PSQL_PASS,default=local_pass" json:"-"`
	Database        string            `env:"PSQL_DBNAME,default=auth_svc"`
	Params          map[string]string `env:"PSQL_PARAMS,default=sslmode:disable" json:",omitempty"`
	MaxOpenConns    int               `env:"PSQL_OPEN_CONNS"`
	MaxIdleConns    int               `env:"PSQL_IDLE_CONNS"`
	ConnMaxLifetime time.Duration     `env:"PSQL_CONN_LIFETIME"`
}

type store struct {
	*sql.DB
	*queries.Queries
}

func New(cfg Config) (Store, error) {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	for k, v := range cfg.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	db, err := sql.Open("pgx", str)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	store := &store{
		DB:      db,
		Queries: queries.New(db),
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
		log.Error("failed to begin transaction", ilog.Err(err))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error("recovered from panic, rolling back transaction and panicking again")
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", ilog.Err(err))
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
