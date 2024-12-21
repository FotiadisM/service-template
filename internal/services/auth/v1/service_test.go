package authv1

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/mocks"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
)

const bufSize = 1024 * 1024

func NewTestServer(t *testing.T, srv authv1.AuthServiceServer) *bufconn.Listener {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, srv)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server exited: %v", err)
		}
	}()

	t.Cleanup(func() {
		grpcServer.GracefulStop()
	})

	return lis
}

func NewTestClient(t *testing.T, srv authv1.AuthServiceServer) authv1.AuthServiceClient {
	t.Helper()

	lis := NewTestServer(t, srv)

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}))
	require.NoError(t, err, "failed to dial connection")

	client := authv1.NewAuthServiceClient(conn)

	t.Cleanup(func() {
		err = conn.Close()
		if err != nil {
			t.Logf("failed to close gRPC client connection %v", err)
		}
	})

	return client
}

type UnitTestingSuite struct {
	grpcServer *grpc.Server
	conn       *grpc.ClientConn

	DB     *mocks.MockDB
	Client authv1.AuthServiceClient
}

func (s *UnitTestingSuite) SetupSuite(t *testing.T) {
	t.Helper()

	s.DB = mocks.NewMockDB(t)
	srv := &Service{db: s.DB}

	lis := bufconn.Listen(bufSize)
	s.grpcServer = grpc.NewServer()
	authv1.RegisterAuthServiceServer(s.grpcServer, srv)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			t.Errorf("gRPC server exited: %v", err)
			os.Exit(1)
		}
	}()

	var err error
	s.conn, err = grpc.NewClient("passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}))
	if err != nil {
		t.Errorf("failed to dial connection: %v", err)
		os.Exit(1)
	}

	s.Client = authv1.NewAuthServiceClient(s.conn)
}

func (s *UnitTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	err := s.conn.Close()
	if err != nil {
		t.Logf("error while closing grpc connection: %v", err)
	}
	s.grpcServer.GracefulStop()
}

func TestUnitTestingSuite(t *testing.T) {
	suite.Run(t, new(UnitTestingSuite))
}
