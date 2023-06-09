package config

import "github.com/knadh/koanf/v2"

type ApiConfig struct {
	Port string
}

func NewApiConfig(k *koanf.Koanf) *ApiConfig {
	return &ApiConfig{
		Port: k.String("api.port"),
	}
}
