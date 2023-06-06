package config

import "github.com/knadh/koanf/v2"

type LoggerConfig struct {
	Name  string
	Level string
}

func NewLoggerConfig(k *koanf.Koanf) *LoggerConfig {
	return &LoggerConfig{
		Name:  k.String("logger.name"),
		Level: k.String("logger.level"),
	}
}
