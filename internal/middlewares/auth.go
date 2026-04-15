package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OnYyon/oregon-api-gateway/internal/clients/sso"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type authPayload struct {
	ginCtx   *gin.Context
	ctx      context.Context
	tokenStr string
	userID   string
	roles    []string
}

type authStep func(*authPayload) error

func AuthMiddleware(client *sso.Client, jwtSecret string, log *slog.Logger) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer("gateway/auth_middleware")
	pipeline := []authStep{
		extractTokenStep(),
		validateSSOStep(client),
		parseClaimsStep(jwtSecret, log),
		enrichContextStep(),
	}

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "Auth.Pipeline", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		reqLog := log.With(slog.String("component", "auth_pipeline"))

		payload := &authPayload{
			ginCtx: c,
			ctx:    ctx,
		}

		for _, step := range pipeline {
			if err := step(payload); err != nil {
				abortPipeline(c, span, reqLog, err)
				return
			}
		}

		c.Request = c.Request.WithContext(payload.ctx)
		c.Next()
	}
}

func abortPipeline(c *gin.Context, span trace.Span, reqLog *slog.Logger, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "authentication_failed")
	reqLog.Warn("authentication pipeline failed", slog.Any("error", err))

	c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
		"error": "unauthorized",
	})
}

func extractTokenStep() authStep {
	return func(p *authPayload) error {
		h := p.ginCtx.GetHeader("Authorization")
		if h == "" {
			return fmt.Errorf("missing_authorization_header")
		}

		if !strings.HasPrefix(h, "Bearer ") {
			return fmt.Errorf("invalid_authorization_format")
		}

		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		if token == "" {
			return fmt.Errorf("empty_token")
		}

		p.tokenStr = token
		return nil
	}
}

func validateSSOStep(client *sso.Client) authStep {
	return func(p *authPayload) error {
		resp, err := client.Validate(p.ctx, &sso.ValidateRequest{
			AccessToken: p.tokenStr,
		})
		if err != nil {
			return fmt.Errorf("sso validation request failed: %w", err)
		}

		if !resp.Validate {
			return fmt.Errorf("token_invalid_by_sso")
		}
		return nil
	}
}

func parseClaimsStep(secret string, reqLog *slog.Logger) authStep {
	return func(p *authPayload) error {
		token, err := jwt.Parse(p.tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			reqLog.Warn("failed to parse jwt or invalid signature", slog.Any("error", err))
			return nil
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil
		}

		if id, ok := claims["id"].(string); ok {
			p.userID = id
		}

		if rolesVal, ok := claims["roles"].([]any); ok {
			for _, r := range rolesVal {
				if strRole, ok := r.(string); ok {
					p.roles = append(p.roles, strRole)
				}
			}
		}

		return nil
	}
}

func enrichContextStep() authStep {
	return func(p *authPayload) error {
		rolesStr := strings.Join(p.roles, ",")

		if p.userID != "" {
			p.ginCtx.Set("user_id", p.userID)
		}
		if len(p.roles) > 0 {
			p.ginCtx.Set("roles", p.roles)
		}

		var mdPairs []string
		if p.userID != "" {
			mdPairs = append(mdPairs, "x-user-id", p.userID)
		}
		if rolesStr != "" {
			mdPairs = append(mdPairs, "x-user-role", rolesStr)
		}

		if len(mdPairs) > 0 {
			p.ctx = metadata.AppendToOutgoingContext(p.ctx, mdPairs...)
		}

		return nil
	}
}
