package grpc

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func Call[T any, R any](
	ctx context.Context,
	conn *grpc.ClientConn,
	log *slog.Logger,
	tracer trace.Tracer,
	timeout time.Duration,
	spanName string,
	fn func(context.Context, *T) (*R, error),
	req *T,
) (*R, error) {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	log.Debug("gRPC call", slog.String("method", spanName))

	resp, err := fn(ctx, req)
	if err != nil {
		log.Error("gRPC call failed",
			slog.String("method", spanName),
			slog.Any("error", err),
			slog.Any("data", req),
		)
		return nil, err
	}

	return resp, nil
}
