package shards

import (
	"fmt"

	"github.com/pijng/mooncache/internal/hasher"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/pijng/mooncache/internal/queue"
)

type cache interface {
	Keymaps() *keymaps.Keymaps
	Queue() *queue.Queue
	PolicyService() *policy.PolicyService
}

type itemArgs interface {
	Key() string
	Value() interface{}
	Cost() int16
	TTL() int64
}

type Shards []Shard
type Shard []interface{}

func Build(amount int8) Shards {
	shards := make(Shards, amount)

	for n := 0; n < int(amount); n++ {
		shards[n] = make(Shard, 0)
	}

	return shards
}

// Set ...
func (s Shards) Set(cluster cache, shardSize int, item itemArgs) error {
	hashedKey := hasher.Sum(item.Key())

	cluster.Queue().Set(hashedKey)
	err := s.set(cluster, shardSize, hashedKey, item)
	cluster.Queue().Release(hashedKey)

	return err
}

func (s Shards) set(cluster cache, shardSize int, hashedKey uint64, item itemArgs) error {
	shardNum := hasher.JCH(hashedKey, int8(len(s)))
	size := lib.ValueSize(item.Value())

	km := cluster.Keymaps()
	ps := cluster.PolicyService()

	if lib.CantFitInShard(km, ps.Variant, shardSize, shardNum, size) {
		return fmt.Errorf("Can't fit value for `%v` key – not enough shard volume: value has `%v` size out of `%v` for shard[%v]",
			item.Key(), size, km.GetShardVolume(shardNum), shardNum)
	}

	ps.EvictUntilCanFit(km, size, shardNum, s.DelByHash)

	lock := km.GetShardLock(shardNum)
	lock.Lock()
	defer lock.Unlock()

	index := len(s[shardNum])
	s[shardNum] = append(s[shardNum], item.Value())

	km.AddKey(hashedKey, index, shardNum, size, item.Cost(), item.TTL())

	return nil
}

// Get ...
func (s Shards) Get(cluster cache, key string) (interface{}, error) {
	hashedKey := hasher.Sum(key)
	if transaction := cluster.Queue().Get(hashedKey); transaction != nil {
		transaction.Wait()
	}

	return s.get(cluster, hashedKey)
}

func (s Shards) get(cluster cache, key uint64) (interface{}, error) {
	shardNum := hasher.JCH(key, int8(len(s)))
	km := cluster.Keymaps()

	lock := km.GetShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := km.GetIndex(key)
	if !ok {
		return nil, lib.ValueNotPresent()
	}

	cluster.PolicyService().UpdateKeyAttrByPolicy(km, key)

	value := s[shardNum][index]
	return value, nil
}

// Del ...
func (s Shards) Del(km *keymaps.Keymaps, key string) {
	s.DelByHash(km, hasher.Sum(key))
}

func (s Shards) DelByHash(km *keymaps.Keymaps, key uint64) {
	shardNum := hasher.JCH(key, int8(len(s)))

	lock := km.GetShardLock(shardNum)
	lock.RLock()
	defer lock.RUnlock()

	index, ok := km.GetIndex(key)
	if !ok {
		return
	}

	s[shardNum][index] = nil
	km.DelKey(key)
}
