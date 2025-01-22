package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/vanguard"
	"github.com/FotiadisM/mock-microservice/internal/config"
	"github.com/stretchr/testify/require"
)

type Server struct {
	server *httptest.Server

	URL    string
	Client *http.Client
}

func NewServer(t *testing.T, config *config.Config, services map[string]http.Handler) *Server {
	t.Helper()

	mux := http.NewServeMux()
	if config.Server.HTTP.DisableRESTTranscoding {
		for path, handler := range services {
			mux.Handle(path, handler)
		}
	} else {
		vanrguardServices := []*vanguard.Service{}
		for path, handler := range services {
			vanrguardServices = append(vanrguardServices, vanguard.NewService(path, handler))
		}
		transcoder, err := vanguard.NewTranscoder(vanrguardServices)
		require.NoError(t, err, "failed to create vanguard transcoder")
		mux.Handle("/", transcoder)
	}

	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()

	return &Server{
		server: server,
		URL:    server.URL,
		Client: server.Client(),
	}
}

func (s *Server) CleanUp() {
	s.server.Close()
}
