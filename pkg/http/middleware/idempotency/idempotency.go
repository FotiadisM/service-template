package idempotency

import (
	"context"
	"errors"
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
	SetKey(ctx context.Context, key string) (bool, error)
	DelKey(ctx context.Context, key string) error
	GetData(ctx context.Context, key string) (*Data, error)
	SetData(ctx context.Context, key string, data *Data, exp time.Duration) error
}

type ErrHandler func(w http.ResponseWriter, r *http.Request, err error) bool

type Middleware struct {
	store         Store
	keyName       string
	replayKeyName string
	dataExp       time.Duration
	errHandler    ErrHandler
}

func NewMiddleware(store Store, errHandler ErrHandler, opts ...Option) *Middleware {
	m := &Middleware{
		store:         store,
		keyName:       "Idempotency-Key",
		replayKeyName: "Idempotent-Replayed",
		dataExp:       3 * time.Minute,
		errHandler:    errHandler,
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
		data, err := m.store.GetData(ctx, key)
		if err == nil {
			maps.Copy(w.Header(), data.Header)
			w.Header().Set(m.replayKeyName, "true")
			w.WriteHeader(data.StatusCode)
			if _, err = w.Write(data.Body); err != nil {
				if m.errHandler(w, r, err) {
					return
				}
			}
			return
		}
		if !errors.Is(err, ErrNoDataFound) {
			if m.errHandler(w, r, err) {
				return
			}
		}

		inProgress, err := m.store.SetKey(ctx, key)
		if err != nil {
			if m.errHandler(w, r, err) {
				return
			}
		}
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
		err = m.store.SetData(ctx, key, data, m.dataExp)
		if err != nil {
			if m.errHandler(w, r, err) {
				return
			}
		}

		err = m.store.DelKey(ctx, key)
		if err != nil {
			if m.errHandler(w, r, err) {
				return
			}
		}
	})
}
