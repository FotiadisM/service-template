package idempotency

import (
	"time"
)

type Option func(*Middleware)

func WithHeaderKeyName(key string) Option {
	return func(m *Middleware) {
		m.keyName = key
	}
}

func WithHeaderReplayKeyName(key string) Option {
	return func(m *Middleware) {
		m.replayKeyName = key
	}
}

func WithDataExp(exp time.Duration) Option {
	return func(m *Middleware) {
		m.dataExp = exp
	}
}
