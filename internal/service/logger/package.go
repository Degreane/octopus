package logger

import (
	"os"
	"strings"
)

var defaultLogger Logger

func init() {
	config := DefaultConfig()

	// Check if we're in a non-interactive environment (like CI/CD)
	// and disable colors if needed
	if os.Getenv("CI") == "true" || os.Getenv("NO_COLOR") == "1" {
		config.UseColors = false
	}

	// Set log level from environment variable if present
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch strings.ToUpper(logLevel) {
		case "DEBUG":
			config.Level = ODebug
		case "INFO":
			config.Level = OInfo
		case "WARN":
			config.Level = OWarn
		case "ERROR":
			config.Level = OError
		case "FATAL":
			config.Level = OFatal
		}
	}

	defaultLogger = New(config)
}

// GetLogger returns the default logger instance
func GetLogger() Logger {
	return defaultLogger
}

// Debug logs a message at debug level
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Info logs a message at info level
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a message at warn level
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs a message at error level
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

// Fatal logs a message at fatal level and exits the program
func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}

// WithField returns a new logger with the given field added to the context
func WithField(key string, value interface{}) Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields returns a new logger with the given fields added to the context
func WithFields(fields map[string]interface{}) Logger {
	return defaultLogger.WithFields(fields)
}
