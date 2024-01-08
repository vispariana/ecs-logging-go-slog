package ecsslog

import (
	"log/slog"
	"runtime/debug"
)

// I'm not quite happy with this approach myself. Keeping it as a reference. 
func Error(err error) slog.Attr {
	msgAttr := slog.String(errorMessageKey, err.Error())
	stAttr := slog.String(errorStackTraceKey, string(debug.Stack()))
	return slog.Group(errorKey, msgAttr, stAttr)
}
