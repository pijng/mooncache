package shards

import (
	"fmt"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/hasher"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/pijng/mooncache/internal/queue"
)

type shard []interface{}

var shards []shard

func Build(amount int8) *[]shard {
	shards = make([]shard, amount)
	for n := 0; n < int(amount); n++ {
		shards[n] = make(shard, 0)
	}

	return &shards
}

// Set ...
func Set(key string, value interface{}, cost int16, ttl int64) error {
	hashedKey := hasher.Sum(key)

	queue.Set(hashedKey)
	err := set(key, hashedKey, value, cost, ttl)
	queue.Release(hashedKey)

	return err
}

func set(key string, hashedKey uint64, value interface{}, cost int16, ttl int64) error {
	shardNum := hasher.JCH(hashedKey, int8(len(shards)))
	size := lib.ValueSize(value)

	if lib.CantFitInShard(config.ShardSize(), shardNum, size) {
		return fmt.Errorf("Can't fit value for `%v` key – not enough shard volume: value has `%v` size out of `%v` for shard[%v]",
			key, size, keymaps.GetShardVolume(shardNum), shardNum)
	}

	policy.EvictUntilCanFit(size, shardNum, DelByHash)

	lock := keymaps.GetShardLock(shardNum)
	lock.Lock()
	defer lock.Unlock()

	pushToShard(shardNum, hashedKey, value, size, cost, ttl)

	return nil
}

func pushToShard(shardNum int, hashedKey uint64, value interface{}, size int, cost int16, ttl int64) {
	index := len(shards[shardNum])
	shards[shardNum] = append(shards[shardNum], value)

	keymaps.AddKey(hashedKey, index, shardNum, size, cost, ttl)
}

// Get ...
func Get(key string) (interface{}, error) {
	hashedKey := hasher.Sum(key)
	if transaction := queue.Get(hashedKey); transaction != nil {
		transaction.Wait()
	}

	return get(hashedKey)
}

func get(key uint64) (interface{}, error) {
	shardNum := hasher.JCH(key, int8(len(shards)))

	lock := keymaps.GetShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := keymaps.GetIndex(key)
	if !ok {
		return nil, lib.ValueNotPresent()
	}

	policy.UpdateKeyAttrByPolicy(key)

	value := shards[shardNum][index]
	return value, nil
}

// Del ...
func Del(key string) {
	DelByHash(hasher.Sum(key))
}

func DelByHash(key uint64) {
	shardNum := hasher.JCH(key, int8(len(shards)))

	lock := keymaps.GetShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := keymaps.GetIndex(key)
	if !ok {
		return
	}

	shards[shardNum][index] = nil
	keymaps.DelKey(key)
}
