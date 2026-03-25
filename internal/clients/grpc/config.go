package grpc

import (
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	Target       string
	Timeout      time.Duration
	DialTimeout  time.Duration
	Interceptors []grpc.UnaryClientInterceptor
}

type Option func(*Config)

func WithTarget(target string) Option {
	return func(c *Config) {
		c.Target = target
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(c *Config) {
		c.DialTimeout = dialTimeout
	}
}

func WithInterceptors(interceptors ...grpc.UnaryClientInterceptor) Option {
	return func(c *Config) {
		c.Interceptors = interceptors
	}
}

func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
