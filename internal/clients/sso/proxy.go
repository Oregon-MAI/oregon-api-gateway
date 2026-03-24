package sso

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func SSOProxy(targetURL string, log *slog.Logger) gin.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Error("invalid SSO target URL", slog.Any("error", err))
		return func(c *gin.Context) {
			c.String(http.StatusServiceUnavailable, "service unavailable")
			c.Abort()
		}
	}

	tracer := otel.GetTracerProvider().Tracer("gateway/sso_proxy")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = nil

	proxy.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetURL(target)
		pr.Out.Host = target.Host

		propagator := otel.GetTextMapPropagator()
		propagator.Inject(pr.Out.Context(), propagation.HeaderCarrier(pr.Out.Header))

		pr.SetXForwarded()
		pr.Out.Header.Del("Connection")
		pr.Out.Header.Del("Upgrade")
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		span := trace.SpanFromContext(r.Context())
		if span.IsRecording() {
			span.RecordError(err)
			span.SetStatus(codes.Error, "proxy_error")
			span.SetAttributes(
				attribute.String("http.path", r.URL.Path),
				attribute.String("http.method", r.Method),
			)
		}

		log.Error("SSO proxy error",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Any("error", err),
		)
		http.Error(w, "SSO service unavailable", http.StatusBadGateway)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		span := trace.SpanFromContext(resp.Request.Context())
		if span.IsRecording() {
			span.SetAttributes(
				attribute.Int("http.status_code", resp.StatusCode),
				attribute.Int64("http.response_content_length", resp.ContentLength),
			)
		}

		resp.Header.Set("X-Proxied-By", "api-gateway")
		return nil
	}

	return func(c *gin.Context) {
		spanName := "SSO.Proxy." + c.Request.Method + c.FullPath()
		ctx, span := tracer.Start(c.Request.Context(), spanName,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", c.Request.URL.Path),
				attribute.String("sso.target", targetURL),
			),
		)

		c.Request = c.Request.WithContext(ctx)
		proxy.ServeHTTP(c.Writer, c.Request)
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		if c.Writer.Status() >= 500 {
			span.SetStatus(codes.Error, "upstream_error")
		}

		span.End()
	}
}
