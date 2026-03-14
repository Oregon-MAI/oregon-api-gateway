package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string       `yaml:"env" env-default:"local"`
	Service string       `yaml:"service"`
	GRPC    GRPCConfig   `yaml:"grpc"`
	Logger  LoggerConfig `yaml:"logger"`
}

type GRPCConfig struct {
	Host string `yaml:"host" env:"HOST"`
	Port int    `yaml:"port" env:"PORT"`
}

type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOGGER_LEVEL"`
	Format string `yaml:"format" env:"LOGGER_FORMAT"`
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
