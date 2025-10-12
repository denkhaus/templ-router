package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/config"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// parseLogLevel converts string to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// createEncoder creates the appropriate encoder based on format
func createEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch strings.ToLower(format) {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "text", "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// createWriteSyncer creates the appropriate write syncer based on config
func createWriteSyncer(cfg *config.Config) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer

	// Always add stdout/stderr based on output config
	switch strings.ToLower(cfg.Logging.Output) {
	case "stdout":
		writers = append(writers, zapcore.AddSync(os.Stdout))
	case "stderr":
		writers = append(writers, zapcore.AddSync(os.Stderr))
	default:
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// Add file output if enabled
	if cfg.Logging.EnableFile {
		// Create directory if it doesn't exist
		dir := filepath.Dir(cfg.Logging.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, zapcore.AddSync(file))
	}

	return zapcore.NewMultiWriteSyncer(writers...), nil
}

// NewService creates the logger service for DI
func NewService(injector do.Injector) (*zap.Logger, error) {
	cfg := do.MustInvoke[*config.Config](injector)

	// Parse log level
	level := parseLogLevel(cfg.Logging.Level)

	// Create encoder based on format
	encoder := createEncoder(cfg.Logging.Format)

	// Create write syncer based on output config
	writeSyncer, err := createWriteSyncer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create write syncer: %w", err)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger with caller info in development
	var logger *zap.Logger
	if cfg.IsDevelopment() {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger = zap.New(core)
	}

	logger.Info("Logger initialized",
		zap.String("mode", func() string {
			if cfg.IsProduction() {
				return "production"
			}
			return "development"
		}()),
		zap.String("level", cfg.Logging.Level),
		zap.String("format", cfg.Logging.Format),
		zap.String("output", cfg.Logging.Output),
		zap.Bool("file_enabled", cfg.Logging.EnableFile),
		zap.String("file_path", func() string {
			if cfg.Logging.EnableFile {
				return cfg.Logging.FilePath
			}
			return "disabled"
		}()),
	)

	return logger, nil
}
