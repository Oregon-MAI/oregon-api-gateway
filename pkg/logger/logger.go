package logger

import "log/slog"

type LoggerConfig struct {
	Level   slog.Level
	Format  string // json or text
	Service string
	Env     string
}

func NewLogger() {

}
