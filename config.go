package mooncache

import (
	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/policy"
)

type Config struct {
	Policy       policy.Policy
	ShardSize    int
	ShardsAmount int8
}

func buildConfig(c *Config) {
	config.Build(c.ShardSize, c.ShardsAmount, c.Policy)
}
