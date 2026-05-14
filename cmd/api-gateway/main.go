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

	bookingclient "github.com/OnYyon/oregon-api-gateway/internal/clients/booking"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/grpc"
	resourceclient "github.com/OnYyon/oregon-api-gateway/internal/clients/resource"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/OnYyon/oregon-api-gateway/internal/config"
	"github.com/OnYyon/oregon-api-gateway/internal/routes"
	"github.com/OnYyon/oregon-api-gateway/pkg/logger"
	"github.com/OnYyon/oregon-api-gateway/pkg/observability/tracer"
	"go.opentelemetry.io/otel"
)

func initTracer(cfg *config.Config, log *slog.Logger) *tracer.Provider {
	tp, err := tracer.New(context.Background(), &tracer.Config{
		ServiceName: cfg.Service,
		EndPoint:    cfg.Trace.EndPoint,
		Insecure:    cfg.Trace.Insecure,
		SampleRatio: cfg.Trace.SampleRatio,
	})
	if err != nil {
		log.Error("failed to init tracer", "error", err)
	}
	return tp
}

func initLogger(cfg *config.Config) (*slog.Logger, *os.File, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err
	}

	f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, nil, err
	}

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

	return log, f, nil
}

func main() {
	cfg := config.MustLoadConfig("./config/local.yml")

	log, f, err := initLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close log file", slog.Any("error", err))
		}
	}()

	tp := initTracer(cfg, log)
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

	resourceClient, err := resourceclient.NewClient(
		grpc.NewConfig(
			grpc.WithTarget(cfg.Resource.PublicTarget),
			grpc.WithTimeout(cfg.Resource.Timeout),
			grpc.WithDialTimeout(cfg.Resource.DialTimeout),
		),
		grpc.NewConfig(
			grpc.WithTarget(cfg.Resource.BookingTarget),
			grpc.WithTimeout(cfg.Resource.Timeout),
			grpc.WithDialTimeout(cfg.Resource.DialTimeout),
		),
		log,
	)
	if err != nil {
		log.Error("failed to create resource client", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := resourceClient.Close(); err != nil {
			log.Error("failed to close resource client", slog.Any("error", err))
		}
	}()

	bookingClient, err := bookingclient.NewClient(
		grpc.NewConfig(
			grpc.WithTarget(cfg.Booking.Target),
			grpc.WithTimeout(cfg.Booking.Timeout),
			grpc.WithDialTimeout(cfg.Booking.DialTimeout),
		),
		log,
	)
	if err != nil {
		log.Error("failed to create booking client", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := bookingClient.Close(); err != nil {
			log.Error("failed to close booking client", slog.Any("error", err))
		}
	}()

	srv := routes.Setup(cfg, log, ssoClient, resourceClient, bookingClient)

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
