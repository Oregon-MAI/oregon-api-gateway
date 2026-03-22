package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.HasTraceID() && spanCtx.HasSpanID() {
		r.AddAttrs(slog.String("trace_id", spanCtx.TraceID().String()))
		r.AddAttrs(slog.String("span_id", spanCtx.SpanID().String()))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{Handler: h.Handler.WithGroup(name)}
}
