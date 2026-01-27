package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func New(level, format string) zerolog.Logger {
	var logger zerolog.Logger

	if format == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	logger = logger.With().Timestamp().Logger()

	switch level {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	return logger
}

type WALogger struct {
	logger zerolog.Logger
	module string
}

func NewWALogger(logger zerolog.Logger, module string) *WALogger {
	return &WALogger{
		logger: logger.With().Str("module", module).Logger(),
		module: module,
	}
}

func (w *WALogger) Debugf(format string, args ...interface{}) {
	w.logger.Debug().Msgf(format, args...)
}

func (w *WALogger) Infof(format string, args ...interface{}) {
	w.logger.Info().Msgf(format, args...)
}

func (w *WALogger) Warnf(format string, args ...interface{}) {
	w.logger.Warn().Msgf(format, args...)
}

func (w *WALogger) Errorf(format string, args ...interface{}) {
	w.logger.Error().Msgf(format, args...)
}

func (w *WALogger) Sub(module string) waLog.Logger {
	return NewWALogger(w.logger, w.module+"/"+module)
}
