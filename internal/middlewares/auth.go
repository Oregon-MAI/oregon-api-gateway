package middlewares

import (
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



func AuthMiddleware(client *sso.Client, jwtSecret string, log *slog.Logger) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer("gateway/auth_middleware")

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "Auth.ValidateToken",
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		reqLog := log.With(slog.String("component", "auth_middleware"))

		tokenStr, err := extractBearerToken(c)
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
			AccessToken: tokenStr,
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

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		var userID string
		var roles []string

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if uuid, ok := claims["id"].(string); ok {
					userID = uuid
				}

				if rolesVal, ok := claims["roles"].([]any); ok {
					for _, r := range rolesVal {
						if strRole, ok := r.(string); ok {
							roles = append(roles, strRole)
						}
					}
				}
			}
		} else {
			reqLog.Warn("failed to parse jwt or invalid signature", slog.Any("error", err))
		}

		rolesStr := strings.Join(roles, ",")

		if userID != "" {
			c.Set("user_id", userID)
			c.Set("roles", roles)
		}

		var mdPairs []string
		if userID != "" {
			mdPairs = append(mdPairs, "x-user-id", userID)
		}
		if rolesStr != "" {
			mdPairs = append(mdPairs, "x-user-role", rolesStr)
		}

		if len(mdPairs) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, mdPairs...)
		}

		c.Request = c.Request.WithContext(ctx)
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
