package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Attributes represents the fields that will be added to every log message
type Attributes struct {
	ServiceName    string
	ServiceVersion string
	CommitSha      string
	BuildTime      string
}

// Config holds the logger attributes and behavior configuration
type Config struct {
	Attributes Attributes
	Level      slog.Level
	AddSource  bool
	Format     Format
}

// NewLogger creates a new logger instance with the provided configuration
// Allows the loglevel to be set from the config
// Allows the AddSource to be set from the config, where this will add the source file and line number to the log output
// Allows the Format to be set from the config (json or text)
// Adds the service name, version, commit sha, and build time to every log message
func NewLogger(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var baseLogger *slog.Logger
	if cfg.Format == FormatJSON {
		baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		baseLogger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	return baseLogger.With(
		"service", cfg.Attributes.ServiceName,
		"version", cfg.Attributes.ServiceVersion,
		"commit_sha", cfg.Attributes.CommitSha,
		"build_time", cfg.Attributes.BuildTime,
	)
}

// ParseLogLevel converts a string log level to a slog.Level
func ParseLogLevel(levelStr string, defaultLevel slog.Level) slog.Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return defaultLevel
	}
}

// ParseFormat converts a string format to a Format type
func ParseFormat(formatStr string, defaultFormat Format) Format {
	switch strings.ToLower(formatStr) {
	case "json":
		return FormatJSON
	case "text":
		return FormatText
	default:
		return defaultFormat
	}
}
