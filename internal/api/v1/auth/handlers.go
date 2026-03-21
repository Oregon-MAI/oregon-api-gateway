package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	ssoClient *sso.Client
	logger    *slog.Logger
	tracer    trace.Tracer
}

func NewHandler(ssoClient *sso.Client, logger *slog.Logger) *Handler {
	return &Handler{
		ssoClient: ssoClient,
		logger:    logger.With(slog.String("component", "auth_handler")),
		tracer:    otel.GetTracerProvider().Tracer("gateway/auth_handler"),
	}
}

func (h *Handler) Login(c *gin.Context) {
	ctx := h.startSpan(c, "Auth.Login")

	var req LoginRequest
	if !h.bind(c, &req) {
		return
	}

	resp, err := h.ssoClient.Login(ctx, &sso.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.handleAuthErr(c, err, "login")
		return
	}

	h.json(c, http.StatusOK, LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *Handler) startSpan(c *gin.Context, name string) context.Context {
	ctx, _ := h.tracer.Start(c.Request.Context(), name,
		trace.WithAttributes(attribute.String("handler", strings.ToLower(strings.TrimPrefix(name, "Auth.")))),
	)
	c.Request = c.Request.WithContext(ctx)
	return ctx
}

func (h *Handler) bind(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		h.json(c, http.StatusBadRequest, errResp("invalid_request", "Invalid request body"))
		return false
	}
	return true
}

func (h *Handler) json(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func (h *Handler) handleAuthErr(c *gin.Context, err error, op string) {
	switch {
	case strings.Contains(err.Error(), "401"), strings.Contains(err.Error(), "unauthorized"):
		h.json(c, http.StatusUnauthorized, errResp("invalid_credentials", "Invalid username or password"))
	case strings.Contains(err.Error(), "unavailable"), strings.Contains(err.Error(), "503"):
		h.json(c, http.StatusServiceUnavailable, errResp("service_unavailable", "Auth service temporarily unavailable"))
	default:
		h.logger.Error("auth operation failed",
			slog.String("op", op),
			slog.Any("error", err),
		)
		h.json(c, http.StatusInternalServerError, errResp("internal_error", "Something went wrong"))
	}
}

func errResp(errType, msg string) map[string]string {
	return map[string]string{"error": errType, "message": msg}
}
