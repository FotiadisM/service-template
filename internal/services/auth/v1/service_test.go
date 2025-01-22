package authv1

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/database"
	"github.com/FotiadisM/mock-microservice/internal/services/auth/v1/queries/mocks"
	"github.com/FotiadisM/mock-microservice/internal/test"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
)

type unitTestingSuiteInternal struct {
	server *test.Server
}

type UnitTestingSuite struct {
	_internal *unitTestingSuiteInternal

	DB      *mocks.MockQuerier
	Service *Service

	ServerURL string
	HTTPClint *http.Client
	Client    authv1connect.AuthServiceClient
}

func (s *UnitTestingSuite) SetupSuite(t *testing.T) {
	t.Helper()

	s.DB = mocks.NewMockQuerier(t)
	s.Service = &Service{db: s.DB}

	config := test.NewConfig()
	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(s.Service,
		connect.WithInterceptors(test.ChainMiddleware(t, config)...),
	)

	server := test.NewServer(t, config, map[string]http.Handler{svcPath: svcHandler})

	s.ServerURL = server.URL
	s.HTTPClint = server.Client
	s.Client = authv1connect.NewAuthServiceClient(server.Client, server.URL)

	s._internal = &unitTestingSuiteInternal{
		server: server,
	}
}

func (s *UnitTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s._internal.server.CleanUp()
}

func TestUnitTestingSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UnitTestingSuite))
}

type endpointTestingSuiteInternal struct {
	postgresContainer *postgres.PostgresContainer
	templateDBName    string
	rootDB            *sql.DB
	server            *test.Server
}

type EndpointTestingSuite struct {
	_internal *endpointTestingSuiteInternal

	DB      *sql.DB
	Service *Service

	ServerURL string
	HTTPClint *http.Client
	Client    authv1connect.AuthServiceClient
}

func (s *EndpointTestingSuite) SetupSuite(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx, "postgres:15.1-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase("test_db"),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "failed to create postgres test container")

	postgresConnURL, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to create postgres connection URL")

	rootDB, err := sql.Open("pgx", postgresConnURL)
	require.NoError(t, err, "failed to open DB connection")

	rootDB.SetMaxOpenConns(1)
	rootDB.SetMaxIdleConns(1)

	templateDBName := "template_db"
	_, err = rootDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", templateDBName))
	require.NoError(t, err, "failed to drop template database")

	_, err = rootDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", templateDBName))
	require.NoError(t, err, "failed to create template database")

	test.ApplyMigrations(ctx, t, strings.ReplaceAll(postgresConnURL, "test_db", templateDBName))
	s.Service = &Service{db: nil}

	config := test.NewConfig()
	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(s.Service,
		connect.WithInterceptors(test.ChainMiddleware(t, config)...),
	)

	server := test.NewServer(t, config, map[string]http.Handler{svcPath: svcHandler})

	s.ServerURL = server.URL
	s.HTTPClint = server.Client
	s.Client = authv1connect.NewAuthServiceClient(server.Client, server.URL)

	s._internal = &endpointTestingSuiteInternal{
		postgresContainer: postgresContainer,
		templateDBName:    templateDBName,
		rootDB:            rootDB,
		server:            server,
	}
}

func (s *EndpointTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s._internal.server.CleanUp()
}

func (s *EndpointTestingSuite) SetupTest(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	testDBName := strings.ReplaceAll(s._internal.templateDBName+"_"+t.Name(), "/", "_")
	testDBName = strings.ToLower(testDBName)

	_, err := s._internal.rootDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	require.NoError(t, err, "failed to drop test database")

	_, err = s._internal.rootDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", testDBName, s._internal.templateDBName))
	require.NoError(t, err, "failed to create test database from template")

	postgresConnURL, err := s._internal.postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to create postgres connection URL")

	conn, err := sql.Open("pgx", strings.ReplaceAll(postgresConnURL, "test_db", testDBName))
	require.NoError(t, err, "failed to open test database")

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	s.DB = conn

	db, err := database.NewFromDBConn(conn)
	require.NoError(t, err, "failed to create db.DB")

	s.Service.db = db
}

func (s *EndpointTestingSuite) TearDownTest(t *testing.T) {
	t.Helper()

	err := s.DB.Close()
	if err != nil {
		t.Logf("failed to close server: %v\n", err)
	}
	err = testcontainers.TerminateContainer(s._internal.postgresContainer)
	if err != nil {
		t.Logf("failed to terminate container: %v\n", err)
	}
}

func TestEndpointTestingSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EndpointTestingSuite))
}
