package logger

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Level       slog.Level
	ServiceName string
	AddSource   bool
	Out         io.Writer
	Format      string // "json" or "text"
	Environment string
}

func New(cfg *Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	out := cfg.Out
	if out == nil {
		out = os.Stdout
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(out, opts)
	default:
		handler = slog.NewTextHandler(out, opts)
	}

	handler = &TraceHandler{Handler: handler}
	logger := slog.New(handler)

	var attrs []any
	if cfg.ServiceName != "" {
		attrs = append(attrs, slog.String("service", cfg.ServiceName))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, slog.String("env", cfg.Environment))
	}

	if len(attrs) > 0 {
		logger = logger.With(attrs...)
	}

	return logger
}
