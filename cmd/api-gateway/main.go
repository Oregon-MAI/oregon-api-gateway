package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/OnYyon/oregon-api-gateway/internal/config"
	"github.com/OnYyon/oregon-api-gateway/internal/routes"
	"github.com/OnYyon/oregon-api-gateway/pkg/logger"
	"github.com/OnYyon/oregon-api-gateway/pkg/observability/tracer"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg := config.MustLoadConfig("./config/local.yml")

	f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close log file", slog.Any("error", err))
		}
	}()

	logCfg := &logger.Config{
		Level:       parseLevel(cfg.Logger.Level),
		Format:      cfg.Logger.Format,
		AddSource:   false,
		Out:         io.MultiWriter(os.Stdout, f),
		ServiceName: cfg.Service,
		Environment: cfg.Env,
	}
	log := logger.New(logCfg)
	slog.SetDefault(log)

	tp, err := tracer.New(context.Background(), &tracer.Config{
		ServiceName: cfg.Service,
		EndPoint:    cfg.Trace.EndPoint,
		Insecure:    cfg.Trace.Insecure,
		SampleRatio: cfg.Trace.SampleRatio,
	})
	if err != nil {
		log.Error("faield to init tracer", "error", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("failed to shutdown tracer", "error", err)
		}
	}()

	ssoClient := sso.NewClient(
		sso.NewConfig(
			sso.WithBaseURL(cfg.SSO.BaseURL),
			sso.WithTimeout(cfg.SSO.Timeout),
		),
		log,
		otel.GetTracerProvider(),
	)

	srv := routes.Setup(cfg, log, ssoClient)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Info("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("shutdown failed", "error", err)
		}
	}()

	log.Info("server starting", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
