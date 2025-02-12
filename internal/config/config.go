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
	Username        string            `env:"USER, required"                  json:"-"`
	Password        string            `env:"PASS, required"                  json:"-"`
	Database        string            `env:"DBNAME, required"`
	Params          map[string]string `env:"PARAMS, default=sslmode:disable"`
	MaxOpenConns    int               `env:"OPEN_CONNS"`
	MaxIdleConns    int               `env:"IDLE_CONNS"`
	ConnMaxLifetime time.Duration     `env:"CONN_LIFETIME"`
}

type Redis struct {
	Host     string `env:"HOST, required"`
	Port     int    `env:"PORT, required"`
	Username string `env:"USER"              json:"-"`
	Password string `env:"PASS"              json:"-"`
	Database int    `env:"DBNAME, default=0"`
}

type Instrumentation struct {
	OtelExporterAddr string `env:"OTEL_EXPORTER_ADDR"`
	OtelSDKDisabled  bool   `env:"OTEL_SDK_DISABLED"`
}

type Cors struct {
	AllowedOrigins      []string `env:"ALLOWED_ORIGINS, default=*"`
	AllowedMethods      []string `env:"ALLOWED_METHODS, default=HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	AllowedHeaders      []string `env:"ALLOWED_HEADERS, default=*"`
	ExposedHeaders      []string `env:"EXPOSED_HEADERS, default=*"`
	MaxAge              int      `env:"MAX_AGE, default=7200"`
	AllowCredentials    bool     `env:"ALLOW_CREDENTIALS, default=false"`
	AllowPrivateNetwork bool     `env:"ALLOW_PRIVATE_NETWORK, default=false"`
}

type Server struct {
	Addr string `env:"ADDR, default=:8080"`

	// Reflection enables gRPC compatible server reflection
	Reflection bool `env:"REFLECTION"`
	// DisableRESTTranscoding disable HTTP JSON+REST to RPC transcoding.
	// It support Google's HTTP transcoding options.
	DisableRESTTranscoding bool `env:"DISABLE_REST_TRANSCODING"`

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

type Config struct {
	Inst    Instrumentation
	Server  Server  `env:", prefix=SERVER_"`
	DB      DB      `env:", prefix=PSQL_"`
	Logging Logging `env:", prefix=LOGGING_"`
	Cors    Cors    `env:", prefix=CORS_"`
	Redis   Redis   `env:", prefix=REDIS_"`
}

func NewConfig(ctx context.Context) *Config {
	config := &Config{}
	if err := envconfig.Process(ctx, config); err != nil {
		fmt.Fprintf(os.Stdout, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	return config
}
