package ilog

type Option func(*handler)

func WithContextKeys(keys ...any) Option {
	return func(h *handler) {
		if len(keys) > 0 {
			h.ctxKeys = keys
		}
	}
}
