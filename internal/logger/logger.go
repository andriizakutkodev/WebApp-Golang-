package logger

import (
	"log/slog"
	"os"
)

var (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

func InitLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case localEnv:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case devEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prodEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
