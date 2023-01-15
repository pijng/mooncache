package mooncache

import (
	"github.com/pijng/mooncache/internal/eviction"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/pijng/mooncache/internal/shards"
)

func buildCache(config *Config) {
	buildConfig(config)
	keymaps.Build(config.ShardsAmount, config.ShardSize)
	shards.Build(config.ShardsAmount)

	if config.Policy != nil {
		policy.Build(config.Policy())
	}

	eviction.Build()
}
