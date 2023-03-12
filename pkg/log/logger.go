package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

const buildTimeFormat = "2006-01-02 15:04:05"

type Logger interface {
	Debug() *zerolog.Event
	Err(error) *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	GetLevel() zerolog.Level
	Hook(zerolog.Hook) zerolog.Logger
	Info() *zerolog.Event
	Level(zerolog.Level) zerolog.Logger
	Log() *zerolog.Event
	Output(io.Writer) zerolog.Logger
	Panic() *zerolog.Event
	Print(...interface{})
	Printf(string, ...interface{})
	Sample(zerolog.Sampler) zerolog.Logger
	Trace() *zerolog.Event
	UpdateContext(func(c zerolog.Context) zerolog.Context)
	Warn() *zerolog.Event
	With() zerolog.Context
	WithContext(context.Context) context.Context
	WithLevel(zerolog.Level) *zerolog.Event
	Write([]byte) (int, error)
}

type LoggerImpl struct {
	zerolog.Logger
}

// NewLogger creates logger that should be used by default.
func NewLogger(opts ...Option) *LoggerImpl {
	o := &options{
		writer: os.Stdout,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.prettify {
		o.writer = zerolog.NewConsoleWriter()
	}

	l := zerolog.New(o.writer).
		Level(ParseLevel(o.level))

	if !o.noTimestamp {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		l = l.Hook(timeHook{})
	}

	c := l.With()

	if o.env != "" {
		c = c.Str("env", o.env)
	}
	if o.buildCommit != "" {
		c = c.Str("build_commit", o.buildCommit)
	}
	if !o.buildTime.IsZero() {
		c = c.Str("build_time", o.buildTime.Format(buildTimeFormat))
	}

	return &LoggerImpl{Logger: c.Logger()}
}

// NewDiscardLogger creates logger that writes nothing.
// Useful in testing and as default value for loggers.
func NewDiscardLogger() *LoggerImpl {
	return &LoggerImpl{
		Logger: zerolog.Nop(),
	}
}

// ParseLevel casts log level string to native zerolog level.
// Useful for getting log levels dynamically without revealing zerolog API:
//
//	var logLevel
//	...
//	logLevel = "warn"
//	...
//	logger.WithLevel(log.ParseLevel(logLevel)).Msg("foo")
func ParseLevel(levelStr string) zerolog.Level {
	level, _ := zerolog.ParseLevel(levelStr)

	//nolint:exhaustive
	switch level {
	default:
		return zerolog.InfoLevel
	case zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
		zerolog.TraceLevel:
	}

	return level
}

// NewNonBlockingWriter writes logs quickly but with no guarantee.
func NewNonBlockingWriter(w io.Writer, size int, pollInterval time.Duration, fallback Logger) io.Writer {
	return diode.NewWriter(w, size, pollInterval, func(missed int) {
		fallback.Warn().Msgf("non-blocking log writer has dropped %v message(s)", missed)
	})
}
