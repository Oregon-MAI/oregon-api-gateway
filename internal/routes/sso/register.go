package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/OnYyon/oregon-api-gateway/internal/api/v1/auth"
	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/OnYyon/oregon-api-gateway/internal/config"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.HTTPConfig, log *slog.Logger, ssoClient *sso.Client) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	authHandler := auth.NewHandler(ssoClient, log)

	pub := r.Group("/api/v1/auth")
	{
		pub.POST("/login", authHandler.Login)
	}

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
