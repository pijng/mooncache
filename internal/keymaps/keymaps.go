package keymaps

import (
	"math"
	"sync"
	"time"
)

type hashmap[K int | uint64, T int | int64] struct {
	Mux *sync.RWMutex
	M   map[K]*T
}

var shardVolumes hashmap[int, int]
var keyIndexes hashmap[uint64, int]
var keyShardNums hashmap[uint64, int]
var valueSizes hashmap[uint64, int]
var valueCosts hashmap[uint64, int]
var valueTTLs hashmap[uint64, int64]
var keyPolicyAttrs hashmap[uint64, int64]

var shardLocks map[int]*sync.RWMutex

func Build(shardsAmount, shardSize int) {
	shardVolumes.Mux = &sync.RWMutex{}
	shardVolumes.M = make(map[int]*int)

	keyIndexes.Mux = &sync.RWMutex{}
	keyIndexes.M = make(map[uint64]*int)

	keyShardNums.Mux = &sync.RWMutex{}
	keyShardNums.M = make(map[uint64]*int)

	valueSizes.Mux = &sync.RWMutex{}
	valueSizes.M = make(map[uint64]*int)

	valueCosts.Mux = &sync.RWMutex{}
	valueCosts.M = make(map[uint64]*int)

	valueTTLs.Mux = &sync.RWMutex{}
	valueTTLs.M = make(map[uint64]*int64)

	keyPolicyAttrs.Mux = &sync.RWMutex{}
	keyPolicyAttrs.M = make(map[uint64]*int64)

	shardLocks = make(map[int]*sync.RWMutex)

	for n := 0; n < shardsAmount; n++ {
		shardVolumes.M[n] = &shardSize
		shardLocks[n] = &sync.RWMutex{}
	}
}

func set[K int | uint64, V int | int64](key K, value V, hm hashmap[K, V]) {
	hm.Mux.Lock()
	defer hm.Mux.Unlock()

	hm.M[key] = &value
}

func get[K int | uint64, V int | int64](key K, hm hashmap[K, V]) (V, bool) {
	hm.Mux.RLock()
	defer hm.Mux.RUnlock()

	v, ok := hm.M[key]
	if !ok {
		return *new(V), false
	}
	return *v, ok
}

func remove[K int | uint64, V int | int64](key K, hm hashmap[K, V]) {
	hm.Mux.Lock()
	defer hm.Mux.Unlock()

	delete(hm.M, key)
}

func AddKey(key uint64, index, shardNum, size, cost int, ttl int64) {
	set(key, index, keyIndexes)
	set(key, cost, valueCosts)
	set(key, ttl, valueTTLs)

	decrementShardVolume(shardNum, size)

	set(key, size, valueSizes)
	set(key, shardNum, keyShardNums)
	set(key, 0, keyPolicyAttrs)
}

func DelKey(key uint64) {
	remove(key, keyIndexes)
	remove(key, valueCosts)
	remove(key, valueTTLs)

	incrementShardVolume(keyShardNum(key), valueSize(key))

	remove(key, valueSizes)
	remove(key, keyShardNums)
	remove(key, keyPolicyAttrs)
}

func decrementShardVolume(shardNum, size int) {
	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(shardNum, shardVolumes)
	set(shardNum, currentVolume-size, shardVolumes)
}

func incrementShardVolume(shardNum, size int) {
	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(shardNum, shardVolumes)
	set(shardNum, currentVolume+size, shardVolumes)
}

func KeyIndex(key uint64) (int, bool) {
	index, ok := get(key, keyIndexes)
	if !ok {
		return 0, false
	}

	return index, true
}

func valueSize(key uint64) int {
	size, ok := get(key, valueSizes)
	if !ok {
		return 0
	}

	return size
}

func valueCost(key uint64) int {
	cost, ok := get(key, valueCosts)
	if !ok {
		return 0
	}

	return cost
}

func ValueTTLs() *hashmap[uint64, int64] {
	return &valueTTLs
}

func StaleKeys() []uint64 {
	now := time.Now().Unix()
	stale := make([]uint64, 0)

	valueTTLs.Mux.Lock()
	defer valueTTLs.Mux.Unlock()

	for key, valueTTL := range valueTTLs.M {
		if *valueTTL > now {
			continue
		}
		stale = append(stale, key)
	}

	return stale
}

func keyShardNum(key uint64) int {
	shardNum, ok := get(key, keyShardNums)
	if !ok {
		return 0
	}

	return shardNum
}

func SetKeyPolicyAttr(key uint64, attr int64) {
	set(key, attr, keyPolicyAttrs)
}

func KeyPolicyAttr(key uint64) (int64, bool) {
	attr, ok := get(key, keyPolicyAttrs)
	if !ok {
		return 0, false
	}

	return attr, true
}

func EnoughSpaceInShard(shardNum, size int) bool {
	volume := ShardVolume(shardNum)
	return volume >= size
}

func ShardVolume(shardNum int) int {
	volume, _ := get(shardNum, shardVolumes)
	return volume
}

func ShardLock(shardNum int) *sync.RWMutex {
	return shardLocks[shardNum]
}

func KeyByMinPolicyAttr() uint64 {
	var hash uint64
	minCost := math.MaxInt
	minValue := int64(math.MaxInt64)

	keyPolicyAttrs.Mux.RLock()
	defer keyPolicyAttrs.Mux.RUnlock()

	for key, attr := range keyPolicyAttrs.M {
		currentAttr := *attr
		currentCost := valueCost(key)

		if currentAttr < minValue && currentCost <= minCost {
			minValue = currentAttr
			minCost = currentCost
			hash = key
		}
	}

	return hash
}

func KeyByMaxPolicyAttr() uint64 {
	var hash uint64
	var maxCost int
	var maxValue int64

	keyPolicyAttrs.Mux.RLock()
	defer keyPolicyAttrs.Mux.RUnlock()

	for key, attr := range keyPolicyAttrs.M {
		currentAttr := *attr
		currentCost := valueCost(key)

		if currentAttr > maxValue && currentCost >= maxCost {
			maxValue = currentAttr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}

func KeyByMinIndex() uint64 {
	var hash uint64
	maxCost := math.MaxInt
	minIndex := int(math.MaxInt64)

	keyPolicyAttrs.Mux.RLock()
	defer keyPolicyAttrs.Mux.RUnlock()

	for key, attr := range keyIndexes.M {
		currentAttr := *attr
		currentCost := valueCost(key)

		if currentAttr < minIndex && currentCost <= maxCost {
			minIndex = currentAttr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}
