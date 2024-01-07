package ecsslog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
)

const (
	ecsVersion = "8.11.0"
	logger     = "log/slog"
)

const (
	ecsVersionKey = "ecs.version"

	timestampKey = "@timestamp"
	messageKey   = "message"
	logLevelKey  = "log.level"
	logLoggerKey = "log.logger"
	fileNameKey  = "file.name"
	fileLineKey  = "file.line"
	logOriginKey = "log.origin"
	functionKey  = "function"

	errorKey           = "error"
	errorMessageKey    = "message"
	errorStackTraceKey = "stack_trace"
)

type Config struct {
	Writer     io.Writer
	LevelNamer func(slog.Level) string
}

func NewHandler(c Config) *Handler {
	if c.LevelNamer == nil {
		c.LevelNamer = defaultLevelNamer
	}
	if c.Writer == nil {
		c.Writer = os.Stdout
	}
	return &Handler{
		next: slog.NewJSONHandler(c.Writer, &slog.HandlerOptions{
			ReplaceAttr: removeJsonHandlerAttrs,
		}),
		levelNamer: c.LevelNamer,
	}
}

func removeJsonHandlerAttrs(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case "time", "msg", "source", "level":
		return slog.Attr{}
	default:
		return a
	}
}

func defaultLevelNamer(l slog.Level) string { return l.String() }

type Handler struct {
	next       slog.Handler
	levelNamer func(slog.Level) string
}

func (x *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return x.next.Enabled(ctx, level)
}

func (x *Handler) Handle(ctx context.Context, record slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{record.PC})
	f, _ := fs.Next()
	record.AddAttrs(
		slog.Time(timestampKey, record.Time),
		slog.String(messageKey, record.Message),
		slog.String(logLevelKey, x.levelNamer(record.Level)),
		slog.String(ecsVersionKey, ecsVersion),
		slog.String(logLoggerKey, logger),
		slog.Group(logOriginKey,
			slog.String(fileNameKey, f.File),
			slog.Int(fileLineKey, f.Line),
			slog.String(functionKey, f.Function),
		),
	)
	return x.next.Handle(ctx, record)
}

func (x *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{next: x.next.WithAttrs(attrs), levelNamer: x.levelNamer}
}

func (x *Handler) WithGroup(name string) slog.Handler {
	return &Handler{next: x.next.WithGroup(name), levelNamer: x.levelNamer}
}
