package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// PostgreSQL databse driver.
	_ "github.com/lib/pq"
)

type Config struct {
	Host            string            `env:"PSQL_HOST,default=localhost"`
	Port            int               `env:"PSQL_PORT,default=5432"`
	Username        string            `env:"PSQL_USER,default=local_user"`
	Password        string            `env:"PSQL_PASS,default=local_pass" json:"-"`
	Database        string            `env:"PSQL_DBNAME,default=local"`
	Params          map[string]string `env:"PSQL_PARAMS,default=sslmode:disable" json:",omitempty"`
	MaxOpenConns    int               `env:"PSQL_OPEN_CONNS"`
	MaxIdleConns    int               `env:"PSQL_IDLE_CONNS"`
	ConnMaxLifetime time.Duration     `env:"PSQL_CONN_LIFETIME"`
}

// ConnectionString generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax.
func ConnectionString(c Config) string {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.Username, c.Password, c.Database)

	for k, v := range c.Params {
		str += fmt.Sprintf(" %s=%s", k, v)
	}

	return str
}

func Open(ctx context.Context, c Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", ConnectionString(c))
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

	return db, nil
}
