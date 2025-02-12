package server

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	bookv1 "github.com/FotiadisM/service-template/api/gen/go/book/v1"
)

type WrappedResponseWriter struct {
	Body       bytes.Buffer
	StatusCode int

	rw http.ResponseWriter
}

func (w *WrappedResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *WrappedResponseWriter) Write(buf []byte) (int, error) {
	return w.Body.Write(buf)
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.rw.WriteHeader(statusCode)
}

func (w *WrappedResponseWriter) Flush() {}

func HTTPTranscoderErrorWrapper(log *slog.Logger, next http.Handler) http.Handler {
	errWriter := connect.NewErrorWriter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if errWriter.IsSupported(r) {
			next.ServeHTTP(w, r)
		}

		wrappedRW := &WrappedResponseWriter{rw: w}
		next.ServeHTTP(wrappedRW, r)

		if wrappedRW.StatusCode == http.StatusOK {
			_, err := w.Write(wrappedRW.Body.Bytes())
			if err != nil {
				log.Error("failed to write to ResponseWriter", "error", err)
			}
			return
		}

		st := &status.Status{}
		err := protojson.Unmarshal(wrappedRW.Body.Bytes(), st)
		if err != nil {
			panic(err)
		}

		if len(st.Details) == 0 {
			return
		}

		m, err := st.Details[0].UnmarshalNew()
		if err != nil {
			log.Error("failed to UnmarshalNew *anypb.Any", "error", err)
			return
		}

		t, ok := m.(*bookv1.Error)
		if ok {
			res := &bookv1.ErrorResponse{
				Error: t,
			}

			var by []byte
			by, err = protojson.Marshal(res)
			if err != nil {
				log.Error("failed to Marshal ErrorResponse", "error", err)
			}
			_, err = w.Write(by)
			if err != nil {
				log.Error("failed to write to ResponseWriter", "error", err)
			}

			return
		}

		log.Warn(fmt.Sprintf("unknown type '%s' passed as error details", st.Details[0].TypeUrl))
		res := &bookv1.ErrorResponse{
			Error: &bookv1.Error{
				Code:        "unknown code",
				Name:        "",
				Description: "",
			},
		}
		by, err := protojson.Marshal(res)
		if err != nil {
			log.Error("failed to Marshal ErrorResponse", "error", err)
		}
		_, err = w.Write(by)
		if err != nil {
			log.Error("failed to write to ResponseWriter", "error", err)
		}

		if len(st.Details) > 1 {
			log.Warn("multiple error details are not supported")
		}
	})
}
