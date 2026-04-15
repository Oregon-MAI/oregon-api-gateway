package grpc

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	log     *slog.Logger
	tracer  trace.Tracer
	timeout time.Duration
}

func NewGRPCClient(cfg *Config, log *slog.Logger) (*Client, error) {
	if cfg.Target == "" {
		return nil, ErrInvalidTarget
	}

	_, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.NewClient(
		cfg.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(cfg.Interceptors...),
	)

	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		log:     log,
		tracer:  otel.GetTracerProvider().Tracer("gateway/grpcClient"),
		timeout: cfg.Timeout,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Timeout() time.Duration {
	return c.timeout
}

func (c *Client) Log() *slog.Logger {
	return c.log
}

func (c *Client) Tracer() trace.Tracer {
	return c.tracer
}
