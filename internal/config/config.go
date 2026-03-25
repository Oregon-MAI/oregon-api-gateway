package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Service  string         `yaml:"service"`
	HTTP     HTTPConfig     `yaml:"http"`
	Logger   LoggerConfig   `yaml:"logger"`
	Trace    TracerConfig   `yaml:"tracer"`
	SSO      SSO            `yaml:"sso"`
	Resource ResourceConfig `yaml:"resource"`
}

type HTTPConfig struct {
	Host         string        `yaml:"host" env:"HTTP_HOST"`
	Port         int           `yaml:"port" env:"HTTP_PORT"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"15s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"15s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type SSO struct {
	BaseURL string        `yaml:"base_url" env:"SSO_BASEURL"`
	Timeout time.Duration `yaml:"timeout" env:"SSO_TIMEOUT"`
}

type ResourceConfig struct {
	BookingTarget string        `yaml:"booking_target" env:"RESOURCE_BOOKING_TARGET"`
	PublicTarget  string        `yaml:"public_target" env:"RESOURCE_PUBLIC_TARGET"`
	Timeout       time.Duration `yaml:"timeout" env:"RESOURCE_TIMEOUT" env-default:"5s"`
	DialTimeout   time.Duration `yaml:"dial_timeout" env:"RESOURCE_DIAL_TIMEOUT" env-default:"5s"`
}

type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOGGER_LEVEL"`
	Format string `yaml:"format" env:"LOGGER_FORMAT"`
}

type TracerConfig struct {
	EndPoint    string  `yaml:"end-point" env:"END_POINT"`
	Insecure    bool    `yaml:"insecure" env:"INSECURE"`
	SampleRatio float64 `yaml:"sample-ratio" env:"SAMPLE_RATION"`
}

func MustLoadConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return cfg
}
