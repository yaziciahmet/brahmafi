package config

import "github.com/knadh/koanf/v2"

type ChainConfig struct {
	RpcUrl    string
	BlockRate int
	Pools     []string
}

func NewChainConfig(k *koanf.Koanf) *ChainConfig {
	return &ChainConfig{
		RpcUrl:    k.String("chain.rpc-url"),
		BlockRate: k.Int("chain.block-rate"),
		Pools:     k.Strings("chain.pools"),
	}
}
