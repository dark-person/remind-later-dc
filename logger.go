package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Create a new formatter for colorized console output.
func ColorFormatter() zerolog.Formatter {
	return func(i any) string {
		level := i.(string)
		switch level {
		case "trace":
			return "\x1b[90m[TRACE]\x1b[0m"
		case "debug":
			return "\x1b[32m[DEBUG]\x1b[0m"
		case "info":
			return "\x1b[34m[INFO ]\x1b[0m"
		case "warn":
			return "\x1b[33m[WARN ]\x1b[0m"
		case "error":
			return "\x1b[31m[ERROR]\x1b[0m"

		default:
			return "[" + level + "]"
		}
	}
}

// Setup global zerolog logger.
func setupLogger(l zerolog.Level) {
	// Console writer
	cw := zerolog.ConsoleWriter{
		Out:         os.Stdout,
		TimeFormat:  "2006-01-02 15:04:05.000",
		FormatLevel: ColorFormatter(),
	}

	// Set up global logger
	mw := zerolog.MultiLevelWriter(cw)
	log.Logger = zerolog.New(mw).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(l)

	// Log something
	log.Debug().Msg("Logger setup completed.")
}
