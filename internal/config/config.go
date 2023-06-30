package config

import (
	"github.com/pijng/mooncache/internal/policy"
)

type Configuration struct {
	Policy       policy.Policy
	ShardSize    int
	ShardsAmount int8
}

var config Configuration

func Build(shardSize int, shardsAmount int8, policy policy.Policy) *Configuration {
	config.ShardSize = shardSize
	config.ShardsAmount = shardsAmount
	config.Policy = policy

	return &config
}

func Config() *Configuration {
	return &config
}

func ShardSize() int {
	return config.ShardSize
}

func ShardsAmount() int8 {
	return config.ShardsAmount
}

func Policy() policy.Policy {
	return config.Policy
}
