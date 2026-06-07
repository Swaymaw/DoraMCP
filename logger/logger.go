package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
)

type colorHandler struct {
	out *os.File
}

func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	color := blue
	label := "DEBUG"

	switch {
	case r.Level >= slog.LevelError:
		color, label = red, "ERROR"
	case r.Level >= slog.LevelWarn:
		color, label = yellow, "WARN"
	case r.Level >= slog.LevelDebug:
		color, label = green, "INFO"
	}

	fmt.Fprintf(h.out, "%s[%s]%s %s %s\n",
		color, label, reset,
		r.Time.Format("2006/01/02 15:04:05"),
		r.Message,
	)
	return nil
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *colorHandler) WithGroup(name string) slog.Handler       { return h }

var Log = slog.New(&colorHandler{out: os.Stdout})
