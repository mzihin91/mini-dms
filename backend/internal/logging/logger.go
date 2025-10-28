package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the structured logger with zerolog
func InitLogger(logLevel string) {
	// Parse log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)

	// Configure pretty console output for development
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Caller().Logger()

	log.Info().
		Str("level", level.String()).
		Msg("Logger initialized successfully")
}

// GetLogger returns a logger with request context
func GetLogger() zerolog.Logger {
	return log.Logger
}
