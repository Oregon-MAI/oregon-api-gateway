package sso

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
	tracer     trace.Tracer
}

func NewClient(cfg *Config, log *slog.Logger, tp trace.TracerProvider) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		logger:  log,
		tracer:  tp.Tracer("gateway/sso_client"),
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: http.DefaultTransport,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, reqBody, respBody any, headers map[string]string, spanName string, attrs ...attribute.KeyValue) error {
	ctx, span := c.tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
	defer span.End()

	var bodyReader io.Reader
	if reqBody != nil {
		body, err := json.Marshal(reqBody)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, bodyReader)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("create request: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed",
			slog.String("method", method),
			slog.String("endpoint", endpoint),
			slog.Any("error", err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("do request %s %s: %w", method, endpoint, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("unexpected status %d from %s: %s", resp.StatusCode, endpoint, string(respData))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if respBody != nil {
		if err := json.Unmarshal(respData, respBody); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}
