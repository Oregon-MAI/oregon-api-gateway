package logging

import (
	"io"
	"log/slog"
)

type LoggerConfig struct {
	Environment string
	ServiceName string
	Format      string // json or text
	AddSource   bool
	Out         io.Writer
	Level       slog.Leveler
}

func NewLogger(cfg *LoggerConfig) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(cfg.Out, &opts)
	} else {
		handler = slog.NewTextHandler(cfg.Out, &opts)
	}

	handler = &TraceIDHandler{handler}
	logger := slog.New(handler).With(
		slog.String("service", cfg.ServiceName),
		slog.String("enviroment", cfg.Environment),
	)
	return logger
}
