package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func AuthMiddleware(client *sso.Client, log *slog.Logger) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer("gateway/auth_middleware")

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "Auth.ValidateToken",
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		reqLog := log.With(slog.String("component", "auth_middleware"))

		token, err := extractBearerToken(c)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "missing_or_invalid_auth_header")
			reqLog.Warn("authentication failed", slog.String("reason", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
			return
		}

		resp, err := client.Validate(ctx, &sso.ValidateRequest{
			AccessToken: token,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "token_validation_failed")
			reqLog.Warn("token validation failed", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid_token",
			})
			return
		}

		if !resp.Validate {
			err := &authError{reason: "token_invalid"}
			span.RecordError(err)
			span.SetStatus(codes.Error, "token_invalid")
			reqLog.Warn("token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid_token",
			})
			return
		}

		c.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, error) {
	h := c.GetHeader("Authorization")
	if h == "" {
		return "", &authError{reason: "missing_authorization_header"}
	}

	if !strings.HasPrefix(h, "Bearer ") {
		return "", &authError{reason: "invalid_authorization_format"}
	}

	token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	if token == "" {
		return "", &authError{reason: "empty_token"}
	}

	return token, nil
}

type authError struct {
	reason string
}

func (e *authError) Error() string {
	return "auth: " + e.reason
}
