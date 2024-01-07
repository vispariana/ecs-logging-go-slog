package ecsslog

import (
	"log/slog"
	"runtime/debug"
)

func Error(err error) slog.Attr {
	msgAttr := slog.String(errorMessageKey, err.Error())
	stAttr := slog.String(errorStackTraceKey, string(debug.Stack()))
	return slog.Group(errorKey, msgAttr, stAttr)
}
