package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/OnYyon/oregon-api-gateway/internal/config"
	logging "github.com/OnYyon/oregon-api-gateway/pkg/logger"
)

func main() {
	cfg := config.MustLoadConfig("./config/local.yml")

	file, err := os.OpenFile("test.log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger := logging.NewLogger(&logging.LoggerConfig{
		Level:       slog.LevelInfo,
		ServiceName: cfg.Service,
		AddSource:   true,
		Out:         file,
		Format:      cfg.Logger.Format,
		Environment: cfg.Env,
	})
	ctx := logging.WithTraceID(context.Background(), "test")
	logger.InfoContext(ctx, "server start on 8000")
}
