package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Logging struct {
	Level     string `env:"LEVEL, default=debug"`
	Output    string `env:"OUTPUT, default=text"`
	AddSource bool   `env:"ADD_SOURCE, default=false"`
}

type DB struct {
	Host            string            `env:"HOST, required"`
	Port            int               `env:"PORT, required"`
	Username        string            `env:"USER, required" json:"-"`
	Password        string            `env:"PASS, required" json:"-"`
	Database        string            `env:"DBNAME, required"`
	Params          map[string]string `env:"PARAMS, default=sslmode:disable" json:",omitempty"`
	MaxOpenConns    int               `env:"OPEN_CONNS"`
	MaxIdleConns    int               `env:"IDLE_CONNS"`
	ConnMaxLifetime time.Duration     `env:"CONN_LIFETIME"`
}

type GRPC struct {
	Addr       string `env:"ADDR, default=:8080"`
	Reflection bool   `env:"REFLECTION"`
	// 	ConnectionTimeout sets the timeout for connection establishment (up to and including HTTP/2 handshaking)
	// for all new connections. A zero or negative value will result in an immediate timeout (120s).
	ConnectionTimeout int64 `env:"CONNECTION_TIMEOUT, default=120"`
	// MaxHeaderListSize sets the max (uncompressed) size of header list that the server is prepared to accept (1MB).
	MaxHeaderListSize uint32 `env:"MAX_HEADER_LIST_SIZE, default=1048576"`
	// MaxRecvMsgSize sets the max message size in bytes the server can receive (4MB).
	MaxRecvMsgSize int `env:"MAX_RECV_MSG_SIZE, default=4194304"`
}

type HTTP struct {
	Addr string `env:"ADDR, default=:9090"`
	// ReadTimeout sets the maximum time a client has to fully stream a request (5s).
	ReadTimeout time.Duration `env:"READ_TIMEOUT, default=5s"`
	// WriteTimeout sets the maximum amount of time a handler has to fully process a request (10s).
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT, default=5s"`
	// IdleTimeout sets the maximum amount of time a Keep-Alive connection can remain idle before
	// being recycled (120s).
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT, default=120s"`
	// ReadHeaderTimeout sets the maximum amount of time a client has to fully stream a request header (5s).
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT, default=$READ_TIMEOUT"`
	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body (1MB).
	MaxHeaderBytes int `env:"MAX_HEADE_RBYTES, default=1048576"`
	// ShutdownTimeout defines how long Graceful shutdown will wait before forcibly shutting down.
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT, default=5s"`
}

type Server struct {
	GRPC GRPC `env:", prefix=GRPC_SERVER_"`
	HTTP HTTP `env:", prefix=HTTPP_SERVER_"`
}

type Config struct {
	Server  Server
	DB      DB      `env:", prefix=PSQL_"`
	Logging Logging `env:", prefix=LOGGING_"`
}

func NewConfig(ctx context.Context) *Config {
	config := &Config{}
	if err := envconfig.Process(ctx, config); err != nil {
		fmt.Fprintf(os.Stdout, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	return config
}
