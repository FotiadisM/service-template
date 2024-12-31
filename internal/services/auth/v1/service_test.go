package authv1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1/authv1connect"
	"github.com/FotiadisM/mock-microservice/internal/db/mocks"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
)

type UnitTestingSuite struct {
	server *httptest.Server

	DB     *mocks.MockDB
	Client authv1connect.AuthServiceClient
}

func (s *UnitTestingSuite) SetupSuite(t *testing.T) {
	t.Helper()

	s.DB = mocks.NewMockDB(t)
	svc := &Service{db: s.DB}

	svcPath, svcHandler := authv1connect.NewAuthServiceHandler(svc,
		connect.WithInterceptors(),
	)

	vanguardSvc := vanguard.NewService(svcPath, svcHandler)
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{vanguardSvc})
	if err != nil {
		t.Errorf("failed to create vanguard transcoder: %v", err)
		t.FailNow()
	}
	mux := http.NewServeMux()
	mux.Handle("/", transcoder)

	s.server = httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	s.Client = authv1connect.NewAuthServiceClient(s.server.Client(), s.server.URL)
}

func (s *UnitTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s.server.Close()
}

func TestUnitTestingSuite(t *testing.T) {
	suite.Run(t, new(UnitTestingSuite))
}
