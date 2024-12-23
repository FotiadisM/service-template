package authv1

import (
	"context"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	authv1 "github.com/FotiadisM/mock-microservice/api/gen/go/auth/v1"
	"github.com/FotiadisM/mock-microservice/internal/db/mocks"
	"github.com/FotiadisM/mock-microservice/pkg/suite"
)

const bufSize = 1024 * 1024

type UnitTestingSuite struct {
	conn       *grpc.ClientConn
	grpcServer *grpc.Server

	DB     *mocks.MockDB
	Client authv1.AuthServiceClient
}

func (s *UnitTestingSuite) SetupSuite(t *testing.T) {
	t.Helper()

	s.DB = mocks.NewMockDB(t)
	srv := &Service{db: s.DB}

	s.grpcServer = grpc.NewServer()
	authv1.RegisterAuthServiceServer(s.grpcServer, srv)

	var err error
	lis := bufconn.Listen(bufSize)
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

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			t.Errorf("failed to start grpc server: %v", err)
			os.Exit(1)
		}
	}()
}

func (s *UnitTestingSuite) TearDownSuite(t *testing.T) {
	t.Helper()

	s.grpcServer.GracefulStop()
	err := s.conn.Close()
	if err != nil {
		t.Logf("error while closing grpc connection: %v", err)
	}
}

func TestUnitTestingSuite(t *testing.T) {
	suite.Run(t, new(UnitTestingSuite))
}
