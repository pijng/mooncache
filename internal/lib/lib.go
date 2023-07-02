package lib

import (
	"fmt"

	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/policy"
)

func ValueSize(value interface{}) int {
	return len([]byte(fmt.Sprint(value)))
}

func CacheNotInitialized() error {
	return fmt.Errorf(`Cache is not initialized, call 'mooncace.New(...)' before calling appropriate methods, preferably during application initialization`)
}

func ValueNotPresent() error {
	return fmt.Errorf("Value  is not present in the cache")
}

func NotEnoughSpace(km *keymaps.Keymaps, variant policy.Variant, shardSize, shardNum int, key string, size int) error {
	enoughSpaceInShard := km.EnoughSpaceInShard(shardNum, size)

	if size > shardSize || !enoughSpaceInShard && variant == "" {
		return fmt.Errorf("Can't fit value for `%v` key – not enough shard volume: value has `%v` size out of `%v` for shard[%v]",
			key, size, km.GetShardVolume(shardNum), shardNum)
	}

	return nil
}
