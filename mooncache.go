package mooncache

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/shards"
)

var once sync.Once

// New ...
func New(config *Config) error {
	if config.ShardsAmount == 0 || config.ShardSize == 0 {
		return fmt.Errorf("ShardsAmount and ShardSize must be greater than 0")
	}

	once.Do(func() {
		buildCache(config)
	})

	return nil
}

// Set ...
func Set(key string, value interface{}, itemOptions ...ItemOptions) error {
	if config.Config() == nil {
		return lib.CacheNotInitialized()
	}

	cost, ttl := getCostTTL(itemOptions)

	if err := shards.Set(key, value, cost, ttl); err != nil {
		return err
	}

	return nil
}

// Get ...
func Get(key string) (interface{}, error) {
	if config.Config() == nil {
		return nil, lib.CacheNotInitialized()
	}

	return shards.Get(key)
}

// Del ...
func Del(key string) error {
	if config.Config() == nil {
		return lib.CacheNotInitialized()
	}

	shards.Del(key)

	return nil
}

func getCostTTL(itemOptions []ItemOptions) (int, int64) {
	var cost int
	ttl := int64(math.MaxInt64)
	if len(itemOptions) > 0 {
		options := itemOptions[0]
		cost = options.Cost
		ttl = time.Now().Add(options.TTL).Unix()
	}

	return cost, ttl
}
