package logger

import (
	"context"
	"os"
	"time"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

type LoggerService struct{}

func NewLoggerService(_ *config.ObservabilityConfig) *LoggerService {
	return &LoggerService{}
}

func (l *LoggerService) Shutdown() {}

func NewLoggerWithService(_ *config.ObservabilityConfig, _ *LoggerService) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(output).With().Timestamp().Logger()
}

func LogBusinessAudit(ctx context.Context, action string, detail string) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Str("action", action).Msg(detail)
}
