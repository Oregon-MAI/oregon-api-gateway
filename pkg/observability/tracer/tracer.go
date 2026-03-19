package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Provider struct {
	tp       *trace.TracerProvider
	shutdown func(context.Context) error
}

func New(ctx context.Context, cfg *Config, opts ...Option) (*Provider, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	for _, opt := range opts {
		opt(cfg)
	}

	exporter, err := createExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := createTracerProvider(exporter, res, cfg.SampleRatio)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{
		tp:       tp,
		shutdown: tp.Shutdown,
	}, nil
}

func (p *Provider) Shutdown(ctx context.Context) error {
	return p.shutdown(ctx)
}
