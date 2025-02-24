package ilog

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func NewLogger(opts ...Option) *slog.Logger {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	replaceAttrFunc := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.LevelKey:
			a.Key = options.LogLevelFieldName
		case slog.TimeKey:
			a.Key = options.TimeFieldName
			a.Value = slog.StringValue(a.Value.Time().Format(options.TimeFieldFormat))
		case slog.MessageKey:
			if options.MessageFieldName != "" {
				a.Key = options.MessageFieldName
			}
		case slog.SourceKey:
			if options.SourceFieldName != "" {
				a.Key = options.SourceFieldName
			}
		}

		if options.ReplaceAttrsOverride != nil {
			return options.ReplaceAttrsOverride(groups, a)
		}
		return a
	}

	handlerOptions := &slog.HandlerOptions{
		Level:       options.LogLevel,
		AddSource:   options.AddSource,
		ReplaceAttr: replaceAttrFunc,
	}

	var logger *slog.Logger
	if options.JSON {
		logger = slog.New(NewFanoutHandler(
			otelslog.NewHandler("book-svc", otelslog.WithSource(options.AddSource)),
			slog.NewJSONHandler(options.Writer, handlerOptions),
		))
	} else {
		logger = slog.New(NewFanoutHandler(
			otelslog.NewHandler("book-svc", otelslog.WithSource(options.AddSource)),
			tint.NewHandler(options.Writer, &tint.Options{
				Level:       options.LogLevel,
				AddSource:   options.AddSource,
				ReplaceAttr: replaceAttrFunc,
			}),
		))
	}

	if len(options.Tags) > 0 {
		group := []any{}
		for k, v := range options.Tags {
			group = append(group, slog.Attr{Key: k, Value: slog.StringValue(v)})
		}
		logger = logger.With(slog.Group("tags", group...))
	}

	return logger
}
