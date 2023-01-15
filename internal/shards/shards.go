package shards

import (
	"fmt"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/hasher"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
)

type shard []interface{}

var shards []shard

func Build(amount int) {
	shards = make([]shard, amount)
	for n := 0; n < amount; n++ {
		shards[n] = make(shard, 0)
	}
}

// Set ...
func Set(key string, value interface{}, cost int, ttl int64) error {
	return set(key, hasher.Sum(key), value, cost, ttl)
}

func set(key string, hashedKey uint64, value interface{}, cost int, ttl int64) error {
	shardNum := hasher.JCH(hashedKey, len(shards))
	size := lib.ValueSize(value)

	if lib.CantFitInShard(config.ShardSize(), shardNum, size) {
		return fmt.Errorf("Can't fit value for `%v` key – not enough shard volume: value has `%v` size out of `%v` for shard[%v]",
			key, size, keymaps.ShardVolume(shardNum), shardNum)
	}

	policy.EvictUntilCanFit(size, shardNum, DelByHash)

	lock := keymaps.ShardLock(shardNum)
	lock.Lock()
	defer lock.Unlock()

	pushToShard(shardNum, hashedKey, value, size, cost, ttl)

	return nil
}

func pushToShard(shardNum int, hashedKey uint64, value interface{}, size, cost int, ttl int64) {
	index := len(shards[shardNum])
	shards[shardNum] = append(shards[shardNum], value)

	keymaps.AddKey(hashedKey, index, shardNum, size, cost, ttl)
}

// Get ...
func Get(key string) (interface{}, error) {
	return get(hasher.Sum(key))
}

func get(key uint64) (interface{}, error) {
	shardNum := hasher.JCH(key, len(shards))

	lock := keymaps.ShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := keymaps.KeyIndex(key)
	if !ok {
		return nil, lib.ValueNotPresent()
	}

	if policy.UpdateKeyAttrByPolicy != nil {
		policy.UpdateKeyAttrByPolicy(key)
	}

	value := shards[shardNum][index]
	return value, nil
}

// Del ...
func Del(key string) {
	DelByHash(hasher.Sum(key))
}

func DelByHash(key uint64) {
	shardNum := hasher.JCH(key, len(shards))

	lock := keymaps.ShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := keymaps.KeyIndex(key)
	if !ok {
		return
	}

	shards[shardNum][index] = nil
	keymaps.DelKey(key)
}
