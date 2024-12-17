package ilog

import (
	"io"
	"log/slog"
	"os"
	"time"
)

type options struct {
	// LogLevel defines the minimum level of severity that app should log.
	// Must be one of:
	// slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError
	LogLevel slog.Level

	// JSON enables structured logging output in json. Should be enabled in production
	JSON bool

	// LogLevelFieldName sets the field name for the log level or severity.
	// Some providers parse and search for different field names.
	// Default is slog.LevelKey
	LogLevelFieldName string

	// MessageFieldName sets the field name for the message.
	// Default is slog.MessageKey
	MessageFieldName string

	// TimeFieldName sets the field name for the time field.
	// Some providers parse and search for different field names.
	// Default is slog.TimeKey
	TimeFieldName string

	// TimeFieldFormat defines the time format of the Time field.
	// Default is time.RFC3339
	TimeFieldFormat string

	// AddSource enables logging the location in the program source code where the logger was called.
	AddSource bool

	// SourceFieldName sets the field anme for the source field.
	// Default is slog.SourceKey
	SourceFieldName string

	// Writer is the log writer, default is os.Stdout
	Writer io.Writer

	// Tags are additional fields included at the root level of all logs.
	// These can be useful for example the commit hash of a build, or an environment
	// name like prod/stg/dev
	Tags map[string]string

	// ReplaceAttrsOverride allows to add custom logic to replace attributes
	// in addition to the default logic set in this package.
	ReplaceAttrsOverride func(groups []string, a slog.Attr) slog.Attr
}

func defaultOptions() *options {
	return &options{
		Writer:            os.Stdout,
		LogLevel:          slog.LevelInfo,
		LogLevelFieldName: slog.LevelKey,
		TimeFieldName:     slog.TimeKey,
		TimeFieldFormat:   time.RFC3339,
		MessageFieldName:  slog.MessageKey,
		SourceFieldName:   slog.SourceKey,
	}
}

type Option func(*options)

func WithLogLevel(lv slog.Level) Option {
	return func(o *options) {
		o.LogLevel = lv
	}
}

func WithJSON(enabled bool) Option {
	return func(o *options) {
		o.JSON = enabled
	}
}

func WithLogLevelFieldName(name string) Option {
	return func(o *options) {
		o.LogLevelFieldName = name
	}
}

func WithMessageFieldName(name string) Option {
	return func(o *options) {
		o.MessageFieldName = name
	}
}

func WithTimeFieldName(name string) Option {
	return func(o *options) {
		o.TimeFieldName = name
	}
}

func WithTimeFieldFormat(format string) Option {
	return func(o *options) {
		o.TimeFieldFormat = format
	}
}

func WithAddSource(enabled bool) Option {
	return func(o *options) {
		o.AddSource = enabled
	}
}

func WithSourceFieldName(name string) Option {
	return func(o *options) {
		o.SourceFieldName = name
	}
}

func WithWrite(w io.Writer) Option {
	return func(o *options) {
		o.Writer = w
	}
}

func WithTags(tags map[string]string) Option {
	return func(o *options) {
		o.Tags = tags
	}
}

func WitReplaceAttrsOverride(fn func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(o *options) {
		o.ReplaceAttrsOverride = fn
	}
}
