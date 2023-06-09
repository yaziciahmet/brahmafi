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
	GetChainConfig() *ChainConfig
}

type KoanfConfiguration struct {
	configEngine *koanf.Koanf
	dbConf       *DatabaseConfig
	loggerConf   *LoggerConfig
	chainConf    *ChainConfig
}

func NewConfiguration() (Configuration, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("common/config/config.yml"), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	config := KoanfConfiguration{
		configEngine: k,
		dbConf:       NewDatabaseConfig(k),
		loggerConf:   NewLoggerConfig(k),
		chainConf:    NewChainConfig(k),
	}

	return &config, nil
}

func (k *KoanfConfiguration) GetDatabaseConfig() *DatabaseConfig {
	return k.dbConf
}

func (k *KoanfConfiguration) GetLoggerConfig() *LoggerConfig {
	return k.loggerConf
}

func (k *KoanfConfiguration) GetChainConfig() *ChainConfig {
	return k.chainConf
}
