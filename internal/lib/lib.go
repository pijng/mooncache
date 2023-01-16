package lib

import (
	"fmt"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/keymaps"
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

func CantFitInShard(shardSize, shardNum int, size int) bool {
	enoughSpaceInShard := keymaps.EnoughSpaceInShard(shardNum, size)
	return size > shardSize || !enoughSpaceInShard && config.Policy() == nil
}
