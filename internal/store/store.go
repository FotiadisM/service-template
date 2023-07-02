package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// PostgreSQL databse driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/FotiadisM/mock-microservice/internal/store/queries"
	"github.com/FotiadisM/mock-microservice/pkg/grpc/interceptor/logger"
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

func New(ctx context.Context, c Config) (Store, error) {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.Username, c.Password, c.Database)
	for k, v := range c.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	db, err := sql.Open("pgx", str)
	if err != nil {
		return nil, err
	}

	if c.MaxOpenConns > 0 {
		db.SetMaxOpenConns(c.MaxOpenConns)
	}
	if c.MaxIdleConns > 0 {
		db.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(c.ConnMaxLifetime)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
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
	log := logger.FromContext(ctx)

	tx, err := s.DB.BeginTx(ctx, options)
	if err != nil {
		log.Error("failed to start transaction", "err", err.Error())
		return err
	}

	defer func() {
		p := recover()
		if p != nil {
			log.Info("recovered from panic, rolling back transaction and panicking again")

			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", "err", err.Error())
			}

			panic(p)
		}
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("failed to roll back transaction", "err", err.Error())
			}
		}
		if err = tx.Commit(); err != nil {
			log.Error("failed to commit transaction", "err", err.Error())
		}
	}()

	err = fn(&store{
		DB:      nil,
		Queries: s.Queries.WithTx(tx),
	})

	return err
}
