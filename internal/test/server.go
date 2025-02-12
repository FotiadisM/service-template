package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FotiadisM/service-template/internal/config"
	"github.com/FotiadisM/service-template/internal/server"
)

type Server struct {
	server *httptest.Server

	URL    string
	Client *http.Client
}

func NewServer(t *testing.T, config *config.Config, services map[string]http.Handler) *Server {
	t.Helper()

	mux := http.NewServeMux()
	if config.Server.DisableRESTTranscoding {
		for path, handler := range services {
			mux.Handle(path, handler)
		}
	} else {
		err := server.HTTPTranscoderHandler(mux, services)
		require.NoError(t, err, "failed to create HTTP transcoder handler")
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
