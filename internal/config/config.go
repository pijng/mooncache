package config

import (
	"github.com/pijng/mooncache/internal/policy"
)

type configuration struct {
	Policy       policy.Policy
	ShardSize    int
	ShardsAmount int
}

var config configuration

func Build(shardSize, shardsAmount int, policy policy.Policy) *configuration {
	config.ShardSize = shardSize
	config.ShardsAmount = shardsAmount
	config.Policy = policy

	return &config
}

func GetConfig() *configuration {
	return &config
}

func GetShardSize() int {
	return config.ShardSize
}

func GetShardsAmount() int {
	return config.ShardsAmount
}

func GetPolicy() policy.Policy {
	return config.Policy
}
