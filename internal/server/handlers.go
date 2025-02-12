package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/grpcreflect"
	"connectrpc.com/vanguard"
	"github.com/rs/cors"

	"github.com/FotiadisM/service-template/api/gen/go/book/v1/bookv1connect"
	"github.com/FotiadisM/service-template/internal/config"
	"github.com/FotiadisM/service-template/pkg/ilog"
)

func CorsHandlers(next http.Handler, config config.Cors) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:      config.AllowedOrigins,
		AllowedMethods:      config.AllowedMethods,
		AllowedHeaders:      config.AllowedHeaders,
		ExposedHeaders:      config.ExposedHeaders,
		MaxAge:              config.MaxAge,
		AllowCredentials:    config.AllowCredentials,
		AllowPrivateNetwork: config.AllowPrivateNetwork,
	})

	return cors.Handler(next)
}

func ReflectionHandler(mux *http.ServeMux, services ...string) {
	reflector := grpcreflect.NewStaticReflector(services...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}

func HTTPTranscoderHandler(mux *http.ServeMux, services map[string]http.Handler) error {
	vanguardServices := []*vanguard.Service{}
	for path, handler := range services {
		vanguardServices = append(vanguardServices,
			vanguard.NewService(path, handler),
		)
	}

	var transcoder *vanguard.Transcoder
	transcoder, err := vanguard.NewTranscoder(vanguardServices)
	if err != nil {
		return fmt.Errorf("failed to create vanguard transcoder: %w", err)
	}
	mux.Handle("/", transcoder)

	return nil
}

func ChainHandlers(
	mux *http.ServeMux,
	config *config.Config,
	log *slog.Logger,
	services map[string]http.Handler,
) http.Handler {
	if config.Server.DisableRESTTranscoding {
		for path, handler := range services {
			mux.Handle(path, handler)
		}
	} else {
		err := HTTPTranscoderHandler(mux, services)
		if err != nil {
			log.Error("failed to http transcoder", ilog.Err(err))
			os.Exit(1)
		}
		log.Info("enabled http rest transcoding")
	}

	if config.Server.Reflection {
		ReflectionHandler(mux, bookv1connect.BookServiceName)
		log.Info("enabled server reflection")
	}

	return CorsHandlers(mux, config.Cors)
}
