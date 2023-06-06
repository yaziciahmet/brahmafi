package config

import "github.com/knadh/koanf/v2"

type DatabaseConfig struct {
	DbUrl         string
	MigrationsDir string
}

func NewDatabaseConfig(k *koanf.Koanf) *DatabaseConfig {
	return &DatabaseConfig{
		DbUrl:         k.String("db.url"),
		MigrationsDir: k.String("db.migrations.dir"),
	}
}
