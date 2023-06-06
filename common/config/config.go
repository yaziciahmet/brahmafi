package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Configuration interface {
	GetDatabaseConfig() *DatabaseConfig
	GetLoggerConfig() *LoggerConfig
}

type KoanfConfig struct {
	configEngine *koanf.Koanf
	dbConf       *DatabaseConfig
	loggerConf   *LoggerConfig
}

func NewConfiguration() (*KoanfConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("common/config/config.yml"), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	config := KoanfConfig{
		configEngine: k,
		dbConf:       NewDatabaseConfig(k),
		loggerConf:   NewLoggerConfig(k),
	}

	return &config, nil
}

func (k *KoanfConfig) GetDatabaseConfig() *DatabaseConfig {
	return k.dbConf
}

func (k *KoanfConfig) GetLoggerConfig() *LoggerConfig {
	return k.loggerConf
}
