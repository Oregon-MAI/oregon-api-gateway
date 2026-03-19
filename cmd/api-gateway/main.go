package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/OnYyon/oregon-api-gateway/internal/config"
	"github.com/OnYyon/oregon-api-gateway/pkg/logger"
	"github.com/OnYyon/oregon-api-gateway/pkg/observability/tracer"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg := config.MustLoadConfig("./config/local.yml")
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log := logger.New(&logger.Config{
		Level:       slog.LevelInfo,
		ServiceName: cfg.Service,
		Out:         logFile,
		AddSource:   true,
		Format:      cfg.Logger.Format,
	})
	tp, err := tracer.New(context.Background(), &tracer.Config{
		ServiceName: cfg.Service,
		SampleRatio: cfg.Trace.SampleRatio,
		EndPoint:    cfg.Trace.EndPoint,
		Insecure:    cfg.Trace.Insecure,
	})
	if err != nil {
		log.Error("Failed to init tracer", "error", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Failed to shutdown tracer", "error", err)
		}
	}()
	tracer := otel.Tracer("api-gw")
	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()
	log.InfoContext(ctx, "Start", slog.Any("user_id", 1))
}
