package middlewares

import (
	"context"
	"errors"
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
	tokenStr string
	userID   string
	roles    []string
	mdPairs  []string
}

type authStep func(context.Context, *authPayload) error

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
		}

		for _, step := range pipeline {
			if err := step(ctx, payload); err != nil {
				abortPipeline(c, span, reqLog, err)
				return
			}
		}

		if len(payload.mdPairs) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, payload.mdPairs...)
		}
		c.Request = c.Request.WithContext(ctx)
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
	return func(_ context.Context, p *authPayload) error {
		h := p.ginCtx.GetHeader("Authorization")
		if h == "" {
			return errors.New("missing_authorization_header")
		}

		if !strings.HasPrefix(h, "Bearer ") {
			return errors.New("invalid_authorization_format")
		}

		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		if token == "" {
			return errors.New("empty_token")
		}

		p.tokenStr = token
		return nil
	}
}

func validateSSOStep(client *sso.Client) authStep {
	return func(ctx context.Context, p *authPayload) error {
		resp, err := client.Validate(ctx, &sso.ValidateRequest{
			AccessToken: p.tokenStr,
		})
		if err != nil {
			return fmt.Errorf("sso validation request failed: %w", err)
		}

		if !resp.Validate {
			return errors.New("token_invalid_by_sso")
		}
		return nil
	}
}

func parseClaimsStep(secret string, reqLog *slog.Logger) authStep {
	return func(_ context.Context, p *authPayload) error {
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

		extractClaims(token, p)
		fmt.Println(p.roles, p.userID)
		return nil
	}
}

func extractClaims(token *jwt.Token, p *authPayload) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
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
}

func enrichContextStep() authStep {
	return func(_ context.Context, p *authPayload) error {
		rolesStr := strings.Join(p.roles, ",")

		if p.userID != "" {
			p.ginCtx.Set("user_id", p.userID)
			p.mdPairs = append(p.mdPairs, "x-user-id", p.userID)
		}

		if len(p.roles) > 0 {
			p.ginCtx.Set("roles", p.roles)
		}

		if rolesStr != "" {
			p.mdPairs = append(p.mdPairs, "x-user-role", rolesStr)
		}

		return nil
	}
}
