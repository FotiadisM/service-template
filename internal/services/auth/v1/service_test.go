package authv1

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"connectrpc.com/vanguard"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/db"
	"github.com/FotiadisM/mock-microservice/internal/services/auth/v1/queries/mocks"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
	"github.com/stretchr/testify/require"
)

type unitTestingSuiteInternal struct {
	server *httptest.Server
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

	validationInterceptor, err := validate.NewInterceptor()
	require.NoError(t, err, "failed to create validation interceptor")

	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(s.Service,
		connect.WithInterceptors(validationInterceptor),
	)

	vanguardSvc := vanguard.NewService(svcPath, svcHandler)
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{vanguardSvc})
	require.NoError(t, err, "failed to create vanguard transcoder")

	mux := http.NewServeMux()
	mux.Handle("/", transcoder)

	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()

	s.ServerURL = server.URL
	s.HTTPClint = server.Client()
	s.Client = authv1connect.NewAuthServiceClient(server.Client(), server.URL)

	s._internal = &unitTestingSuiteInternal{
		server: server,
	}
}

func (s *UnitTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s._internal.server.Close()
}

func TestUnitTestingSuite(t *testing.T) {
	suite.Run(t, new(UnitTestingSuite))
}

type endpointTestingSuiteInternal struct {
	templateDBName string
	rootDB         *sql.DB
	server         *httptest.Server
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
	rootDB, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres")
	require.NoError(t, err, "failed to open DB connection")

	rootDB.SetMaxOpenConns(1)
	rootDB.SetMaxIdleConns(1)

	templateDBName := "template_db"
	_, err = rootDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", templateDBName))
	require.NoError(t, err, "failed to drop template database")

	_, err = rootDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", templateDBName))
	require.NoError(t, err, "failed to create template database")

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("../../../db/migrations/"),
		),
	)
	require.NoError(t, err, "failed to create atlas workdir")
	defer workdir.Close()

	atlasClient, err := atlasexec.NewClient(workdir.Path(), "atlas")
	require.NoError(t, err, "failed to create atlas client")

	res, err := atlasClient.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", templateDBName),
	})
	require.NoError(t, err, "failed to apply migrations")
	t.Logf("Applied %d migrations\n", len(res.Applied))

	s.Service = &Service{db: nil}

	validationInterceptor, err := validate.NewInterceptor()
	require.NoError(t, err, "failed to create validation interceptor")

	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(s.Service,
		connect.WithInterceptors(validationInterceptor),
	)

	vanguardSvc := vanguard.NewService(svcPath, svcHandler)
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{vanguardSvc})
	require.NoError(t, err, "failed to create vanguard transcoder")

	mux := http.NewServeMux()
	mux.Handle("/", transcoder)

	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()

	s.ServerURL = server.URL
	s.HTTPClint = server.Client()
	s.Client = authv1connect.NewAuthServiceClient(server.Client(), server.URL)

	s._internal = &endpointTestingSuiteInternal{
		templateDBName: templateDBName,
		rootDB:         rootDB,
		server:         server,
	}
}

func (s *EndpointTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s._internal.server.Close()
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

	conn, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname="+testDBName)
	require.NoError(t, err, "failed to open test database")

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	s.DB = conn

	db, err := db.NewFromDBConn(conn)
	require.NoError(t, err, "failed to create db.DB")

	s.Service.db = db
}

func (s *EndpointTestingSuite) TearDownTest(t *testing.T) {
	t.Helper()

	_ = s.DB.Close()
}

func TestEndpointTestingSuite(t *testing.T) {
	suite.Run(t, new(EndpointTestingSuite))
}
