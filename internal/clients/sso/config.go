package sso

import "time"

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type Option func(*Config)

func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
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
