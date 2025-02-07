package idempotency

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"time"

	"github.com/FotiadisM/mock-microservice/pkg/http/middleware"
)

var ErrNoDataFound = errors.New("no data found")

type Data struct {
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
	StatusCode int         `json:"statusCode"`
}

type Store interface {
	SetKey(ctx context.Context, key string) bool
	DelKey(ctx context.Context, key string)
	GetData(ctx context.Context, key string) *Data
	SetData(ctx context.Context, key string, data *Data, exp time.Duration)
}

type Middleware struct {
	store         Store
	keyName       string
	replayKeyName string
	dataExp       time.Duration
	log           *slog.Logger
}

func NewMiddleware(store Store, opts ...Option) *Middleware {
	m := &Middleware{
		store:         store,
		keyName:       "Idempotency-Key",
		replayKeyName: "Idempotent-Replayed",
		dataExp:       3 * time.Minute,
		log:           slog.Default(),
	}
	for _, o := range opts {
		o(m)
	}

	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(m.keyName)
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		data := m.store.GetData(ctx, key)
		if data != nil {
			maps.Copy(w.Header(), data.Header)
			w.Header().Set(m.replayKeyName, "true")
			w.WriteHeader(data.StatusCode)
			_, err := w.Write(data.Body)
			if err != nil {
				m.log.Error("http-idempotency: failed to write response to ResponseWriter", "error", err)
			}

			return
		}

		inProgress := m.store.SetKey(ctx, key)
		if inProgress {
			w.WriteHeader(http.StatusConflict)
		}

		wrappedRW := &middleware.WrappedResponseWriter{ResponseWriter: w}
		next.ServeHTTP(wrappedRW, r)

		data = &Data{
			Header:     w.Header().Clone(),
			Body:       wrappedRW.Body.Bytes(),
			StatusCode: wrappedRW.StatusCode,
		}
		if data.StatusCode == 0 {
			data.StatusCode = http.StatusOK
		}

		m.store.SetData(ctx, key, data, m.dataExp)
		m.store.DelKey(ctx, key)
	})
}
