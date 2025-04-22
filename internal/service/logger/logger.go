package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of the log message
type LogLevel int

const (
	// Debug level for detailed information
	ODebug LogLevel = iota
	// Info level for general operational information
	OInfo
	// Warn level for non-critical issues
	OWarn
	// Error level for errors that should be addressed
	OError
	// Fatal level for critical errors that require immediate attention
	OFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case ODebug:
		return "DEBUG"
	case OInfo:
		return "INFO"
	case OWarn:
		return "WARN"
	case OError:
		return "ERROR"
	case OFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color code for the log level
func (l LogLevel) Color() string {
	switch l {
	case ODebug:
		return "\033[36m" // Cyan
	case OInfo:
		return "\033[32m" // Green
	case OWarn:
		return "\033[33m" // Yellow
	case OError:
		return "\033[31m" // Red
	case OFatal:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Default
	}
}

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// OctopusLogger implements the Logger interface
type OctopusLogger struct {
	level     LogLevel
	output    io.Writer
	fields    map[string]interface{}
	useColors bool
}

// Config holds the configuration for the logger
type Config struct {
	Level     LogLevel
	Output    io.Writer
	UseColors bool
}

// DefaultConfig returns a default configuration for the logger
func DefaultConfig() Config {
	return Config{
		Level:     OInfo,
		Output:    os.Stdout,
		UseColors: true,
	}
}

// New creates a new logger with the given configuration
func New(config Config) Logger {
	return &OctopusLogger{
		level:     config.Level,
		output:    config.Output,
		fields:    make(map[string]interface{}),
		useColors: config.UseColors,
	}
}

// WithField returns a new logger with the given field added to the context
func (l *OctopusLogger) WithField(key string, value interface{}) Logger {
	newLogger := &OctopusLogger{
		level:     l.level,
		output:    l.output,
		fields:    make(map[string]interface{}),
		useColors: l.useColors,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value
	return newLogger
}

// WithFields returns a new logger with the given fields added to the context
func (l *OctopusLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := &OctopusLogger{
		level:     l.level,
		output:    l.output,
		fields:    make(map[string]interface{}),
		useColors: l.useColors,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// Debug logs a message at debug level
func (l *OctopusLogger) Debug(msg string, args ...interface{}) {
	if l.level <= ODebug {
		l.log(ODebug, msg, args...)
	}
}

// Info logs a message at info level
func (l *OctopusLogger) Info(msg string, args ...interface{}) {
	if l.level <= OInfo {
		l.log(OInfo, msg, args...)
	}
}

// Warn logs a message at warn level
func (l *OctopusLogger) Warn(msg string, args ...interface{}) {
	if l.level <= OWarn {
		l.log(OWarn, msg, args...)
	}
}

// Error logs a message at error level
func (l *OctopusLogger) Error(msg string, args ...interface{}) {
	if l.level <= OError {
		l.log(OError, msg, args...)
	}
}

// Fatal logs a message at fatal level and exits the program
func (l *OctopusLogger) Fatal(msg string, args ...interface{}) {
	if l.level <= OFatal {
		l.log(OFatal, msg, args...)
		os.Exit(1)
	}
}

// log formats and writes the log message
func (l *OctopusLogger) log(level LogLevel, msg string, args ...interface{}) {
	// Format the message if there are arguments
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// Extract just the filename from the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Start building the log entry
	var builder strings.Builder

	// Add color if enabled
	if l.useColors {
		builder.WriteString(level.Color())
	}

	// Write timestamp, level, and file information
	builder.WriteString(fmt.Sprintf("[%s] [%s] [%s:%d] ", timestamp, level.String(), file, line))

	// Write message
	builder.WriteString(msg)

	// Write fields if any
	if len(l.fields) > 0 {
		builder.WriteString(" {")
		first := true
		for k, v := range l.fields {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s=%v", k, v))
			first = false
		}
		builder.WriteString("}")
	}

	// Reset color if enabled
	if l.useColors {
		builder.WriteString("\033[0m")
	}

	// Add newline
	builder.WriteString("\n")

	// Write to output
	fmt.Fprint(l.output, builder.String())
}
