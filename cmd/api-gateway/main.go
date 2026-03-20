package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/OnYyon/oregon-api-gateway/pkg/logger"
	"github.com/OnYyon/oregon-api-gateway/pkg/observability/tracer"
	"github.com/ilyakaznacheev/cleanenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Config struct {
	Env     string `yaml:"env" env-default:"local"`
	Service string `yaml:"service"`

	Logger struct {
		Level     slog.Level `yaml:"level"`
		Format    string     `yaml:"format"`
		FilePath  string     `yaml:"file_path"`
		AddSource bool       `yaml:"add_source"`
	} `yaml:"logger"`

	Tracer tracer.Config `yaml:"tracer"`
}

func main() {
	var cfg Config
	configPath := "config/local.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "./config/local.yml"
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	var logOutput io.Writer = os.Stdout
	if cfg.Logger.FilePath != "" {
		f, err := os.OpenFile(cfg.Logger.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			panic(fmt.Errorf("failed to open log file: %w", err))
		}
		defer func() {
			if err := f.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to close log file: %v\n", err)
			}
		}()
		logOutput = io.MultiWriter(os.Stdout, f)
	}

	logCfg := &logger.Config{
		Level:       cfg.Logger.Level,
		Format:      cfg.Logger.Format,
		AddSource:   cfg.Logger.AddSource,
		Out:         logOutput,
		ServiceName: cfg.Service,
		Environment: cfg.Env,
	}
	log := logger.New(logCfg)
	slog.SetDefault(log)

	log.Info("Starting application simulation", "config", cfg)
	if cfg.Tracer.ServiceName == "" {
		cfg.Tracer.ServiceName = cfg.Service
	}

	tp, err := tracer.New(context.Background(), &cfg.Tracer)
	if err != nil {
		log.Error("Failed to init tracer", "error", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Failed to shutdown tracer", "error", err)
		}
	}()

	ctx := context.Background()
	handler := &UserHandler{
		svc: &UserService{
			repo: &UserRepository{},
		},
	}

	fmt.Println("\n--- Scenario 1: Successful Request ---")
	handler.HandleCreateUser(ctx, "john_doe")
	fmt.Println("\n--- Scenario 2: Failed Request (DB Error) ---")
	handler.HandleCreateUser(ctx, "error_user")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\nSimulation complete. Check app.log for logs with consistent trace_ids.")
}

type UserHandler struct {
	svc *UserService
}

func (h *UserHandler) HandleCreateUser(ctx context.Context, username string) {
	tr := otel.Tracer("http-handler")
	ctx, span := tr.Start(ctx, "POST /users")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", "POST"),
		attribute.String("http.route", "/users"),
		attribute.String("client.ip", "192.168.1.10"),
	)

	slog.InfoContext(ctx, "Handling create user request", "username", username)

	err := h.svc.CreateUser(ctx, username)
	if err != nil {
		slog.ErrorContext(ctx, "Request failed", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	slog.InfoContext(ctx, "Request processed successfully", "username", username)
	span.SetStatus(codes.Ok, "OK")
}

type UserService struct {
	repo *UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, username string) error {
	tr := otel.Tracer("user-service")
	ctx, span := tr.Start(ctx, "UserService.CreateUser")
	defer span.End()
	slog.InfoContext(ctx, "Validating user", "username", username)
	// #nosec G404 - Not security-sensitive, just for sleep jitter
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	if username == "" {
		return errors.New("username cannot be empty")
	}

	err := s.repo.Save(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

type UserRepository struct{}

func (r *UserRepository) Save(ctx context.Context, username string) error {
	tr := otel.Tracer("user-repository")
	ctx, span := tr.Start(ctx, "UserRepository.Save")
	defer span.End()
	span.SetAttributes(attribute.String("db.system", "postgres"))
	slog.InfoContext(ctx, "Saving user to database", "username", username)
	// #nosec G404 - Not security-sensitive, just for sleep jitter
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	if username == "error_user" {
		err := errors.New("database connection timeout")
		slog.ErrorContext(ctx, "Database error", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	slog.InfoContext(ctx, "User saved successfully")
	return nil
}
