package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/netprobe/internal/config"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

type logrusLogger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

func New(cfg config.LoggingConfig) (Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// Set formatter
	switch cfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	// Set output
	if cfg.File != "" {
		dir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Write to both file and stdout
		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(multiWriter)
	}

	return &logrusLogger{
		logger: logger,
		fields: make(logrus.Fields),
	}, nil
}

func (l *logrusLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(l.fields).WithFields(l.parseFields(keysAndValues...)).Debug(msg)
}

func (l *logrusLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(l.fields).WithFields(l.parseFields(keysAndValues...)).Info(msg)
}

func (l *logrusLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(l.fields).WithFields(l.parseFields(keysAndValues...)).Warn(msg)
}

func (l *logrusLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(l.fields).WithFields(l.parseFields(keysAndValues...)).Error(msg)
}

func (l *logrusLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.WithFields(l.fields).WithFields(l.parseFields(keysAndValues...)).Fatal(msg)
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &logrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func (l *logrusLogger) parseFields(keysAndValues ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if ok {
				fields[key] = keysAndValues[i+1]
			}
		}
	}
	return fields
}

// Test helpers
func NewTestLogger(t interface{}) Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &logrusLogger{logger: logger, fields: make(logrus.Fields)}
}

func NewNullLogger() Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return &logrusLogger{logger: logger, fields: make(logrus.Fields)}
}
