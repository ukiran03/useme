package clog

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var logger *slog.Logger

func init() {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindAny {
				// writes all errors in red
				if _, ok := a.Value.Any().(error); ok {
					return tint.Attr(9, a)
				}
			}
			return a
		},
	})

	logger = slog.New(handler)

	// Set as global slog default so standard slog calls use our theme
	slog.SetDefault(logger)
}

func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }
