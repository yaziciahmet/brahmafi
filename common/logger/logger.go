package logger

import (
	"brahmafi/common/config"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Clone(name string) Logger
	GetWriter() io.Writer
}

type hcLogger struct {
	Logger

	logger hclog.Logger
	config *config.LoggerConfig
}

func NewLogger(config *config.LoggerConfig) Logger {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:                 config.Name,
		Level:                hclog.LevelFromString(config.Level),
		Color:                hclog.AutoColor,
		ColorHeaderAndFields: true,
	})

	return &hcLogger{
		logger: logger,
		config: config,
	}
}

func (h *hcLogger) Trace(msg string, args ...interface{}) {
	h.logger.Trace(msg, args...)
}

func (h *hcLogger) Debug(msg string, args ...interface{}) {
	h.logger.Debug(msg, args...)
}

func (h *hcLogger) Info(msg string, args ...interface{}) {
	h.logger.Info(msg, args...)
}

func (h *hcLogger) Warn(msg string, args ...interface{}) {
	h.logger.Warn(msg, args...)
}

func (h *hcLogger) Error(msg string, args ...interface{}) {
	h.logger.Error(msg, args...)
}

func (h *hcLogger) Fatal(msg string, args ...interface{}) {
	h.logger.Error(msg, args...)
	os.Exit(1)
}

func (h *hcLogger) Clone(name string) Logger {
	return &hcLogger{logger: h.logger.Named(name)}
}

func (h *hcLogger) GetWriter() io.Writer {
	return h.logger.StandardWriter(&hclog.StandardLoggerOptions{})
}
