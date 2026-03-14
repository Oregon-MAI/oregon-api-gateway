package logging

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

var traceIDKey ctxKey

type TraceIDHandler struct {
	slog.Handler
}

func (h *TraceIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *TraceIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceIDHandler{h.Handler.WithAttrs(attrs)}
}

func (h *TraceIDHandler) WithGroup(name string) slog.Handler {
	return &TraceIDHandler{h.Handler.WithGroup(name)}
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}
