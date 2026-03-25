package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		spanName := c.Request.Method + " " + c.FullPath()
		if spanName == "GET " || spanName == "POST " {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		attrs := []attribute.KeyValue{
			semconv.HTTPRequestMethodOriginal(c.Request.Method),
			semconv.URLFull(c.Request.URL.String()),
			semconv.ClientAddress(c.ClientIP()),
			semconv.URLScheme(c.Request.URL.Scheme),
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		propagator.Inject(ctx, propagation.HeaderCarrier(c.Writer.Header()))
		if sc := span.SpanContext(); sc.IsValid() {
			c.Header("X-Trace-ID", sc.TraceID().String())
		}

		c.Next()
		span.SetAttributes(semconv.HTTPResponseStatusCode(c.Writer.Status()))
		if c.Writer.Status() >= 500 {
			span.SetStatus(codes.Error, "Internal Server Error")
		}
	}
}
