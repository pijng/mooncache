package mooncache

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/pijng/mooncache/internal/eviction"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/pijng/mooncache/internal/queue"
	"github.com/pijng/mooncache/internal/shards"
)

type ItemOptions struct {
	TTL  time.Duration
	Cost int16
}

type cache struct {
	policyService *policy.PolicyService
	config        *Config
	keymaps       *keymaps.Keymaps
	queue         *queue.Queue
	shards        shards.Shards
}

func (c *cache) Config() *Config                      { return c.config }
func (c *cache) Keymaps() *keymaps.Keymaps            { return c.keymaps }
func (c *cache) Shards() shards.Shards                { return c.shards }
func (c *cache) Queue() *queue.Queue                  { return c.queue }
func (c *cache) PolicyService() *policy.PolicyService { return c.policyService }

type itemArgs struct {
	value interface{}
	key   string
	ttl   int64
	cost  int16
}

func (i *itemArgs) Key() string        { return i.key }
func (i *itemArgs) Value() interface{} { return i.value }
func (i *itemArgs) Cost() int16        { return i.cost }
func (i *itemArgs) TTL() int64         { return i.ttl }

var cluster *cache
var once sync.Once

// New ...
func New(config *Config) error {
	if config.ShardsAmount == 0 || config.ShardSize == 0 {
		return fmt.Errorf("ShardsAmount and ShardSize must be greater than 0")
	}

	once.Do(func() {
		km := keymaps.Build(config.ShardsAmount, config.ShardSize)
		s := shards.Build(config.ShardsAmount)
		q := queue.Build()

		cluster = &cache{
			config:  config,
			keymaps: km,
			shards:  s,
			queue:   q,
		}

		if config.Algorithm() != "" {
			policyService := policy.Build(cluster.Keymaps(), config.Algorithm())
			cluster.policyService = policyService
		}

		eviction.Run(cluster.Shards(), cluster.Keymaps())
	})

	return nil
}

func (c *cache) Set(key string, value interface{}, itemOptions ...ItemOptions) error {
	return Set(key, value, itemOptions...)
}

func (c *cache) Get(key string) (interface{}, error) {
	return Get(key)
}

func (c *cache) Del(key string) error {
	return Del(key)
}

// Set ...
func Set(key string, value interface{}, itemOptions ...ItemOptions) error {
	if cluster == nil {
		return lib.CacheNotInitialized()
	}

	cost, ttl := getCostTTL(itemOptions)

	item := &itemArgs{
		key:   key,
		value: value,
		cost:  cost,
		ttl:   ttl,
	}

	if err := cluster.Shards().Set(cluster, cluster.Config().ShardSize, item); err != nil {
		return err
	}

	return nil
}

// Get ...
func Get(key string) (interface{}, error) {
	if cluster == nil {
		return nil, lib.CacheNotInitialized()
	}

	return cluster.Shards().Get(cluster, key)
}

// Del ...
func Del(key string) error {
	if cluster == nil {
		return lib.CacheNotInitialized()
	}

	cluster.Shards().Del(cluster.Keymaps(), key)

	return nil
}

func getCostTTL(itemOptions []ItemOptions) (int16, int64) {
	var cost int16
	ttl := int64(math.MaxInt64)

	if len(itemOptions) > 0 {
		options := itemOptions[0]
		cost = options.Cost
		ttl = time.Now().Add(options.TTL).Unix()
	}

	return cost, ttl
}
