package keymaps

import (
	"math"
	"sync"
	"time"
)

type hashmap[K int | uint64, T int | int16 | int64] struct {
	Mux *sync.RWMutex
	M   map[K]T
}

var shardVolumes hashmap[int, int]
var indexes hashmap[uint64, int]
var shardNums hashmap[uint64, int]
var valueSizes hashmap[uint64, int]
var valueCosts hashmap[uint64, int16]
var valueTTLs hashmap[uint64, int64]
var policyAttrs hashmap[uint64, int64]

var shardLocks map[int]*sync.RWMutex

func Build(shardsAmount int8, shardSize int) {
	shardVolumes.Mux = &sync.RWMutex{}
	shardVolumes.M = make(map[int]int)

	indexes.Mux = &sync.RWMutex{}
	indexes.M = make(map[uint64]int)

	shardNums.Mux = &sync.RWMutex{}
	shardNums.M = make(map[uint64]int)

	valueSizes.Mux = &sync.RWMutex{}
	valueSizes.M = make(map[uint64]int)

	valueCosts.Mux = &sync.RWMutex{}
	valueCosts.M = make(map[uint64]int16)

	valueTTLs.Mux = &sync.RWMutex{}
	valueTTLs.M = make(map[uint64]int64)

	policyAttrs.Mux = &sync.RWMutex{}
	policyAttrs.M = make(map[uint64]int64)

	shardLocks = make(map[int]*sync.RWMutex)

	for n := 0; n < int(shardsAmount); n++ {
		shardVolumes.M[n] = shardSize
		shardLocks[n] = &sync.RWMutex{}
	}
}

func set[K int | uint64, V int | int16 | int64](hm hashmap[K, V], key K, value V) {
	hm.Mux.Lock()
	defer hm.Mux.Unlock()

	hm.M[key] = value
}

func get[K int | uint64, V int | int16 | int64](hm hashmap[K, V], key K) (V, bool) {
	hm.Mux.RLock()
	defer hm.Mux.RUnlock()

	v, ok := hm.M[key]
	if !ok {
		return *new(V), false
	}
	return v, ok
}

func remove[K int | uint64, V int | int16 | int64](hm hashmap[K, V], key K) {
	hm.Mux.Lock()
	defer hm.Mux.Unlock()

	delete(hm.M, key)
}

func AddKey(key uint64, index, shardNum, size int, cost int16, ttl int64) {
	set(indexes, key, index)
	set(valueCosts, key, cost)
	set(valueTTLs, key, ttl)

	decrementShardVolume(shardNum, size)

	set(valueSizes, key, size)
	set(shardNums, key, shardNum)
	set(policyAttrs, key, 0)
}

func DelKey(key uint64) {
	remove(indexes, key)
	remove(valueCosts, key)
	remove(valueTTLs, key)

	incrementShardVolume(getShardNum(key), getValueSize(key))

	remove(valueSizes, key)
	remove(shardNums, key)
	remove(policyAttrs, key)
}

func decrementShardVolume(shardNum, size int) {
	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(shardVolumes, shardNum)
	set(shardVolumes, shardNum, currentVolume-size)
}

func incrementShardVolume(shardNum, size int) {
	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(shardVolumes, shardNum)
	set(shardVolumes, shardNum, currentVolume+size)
}

func GetIndex(key uint64) (int, bool) {
	index, ok := get(indexes, key)
	if !ok {
		return 0, false
	}

	return index, true
}

func getValueSize(key uint64) int {
	size, ok := get(valueSizes, key)
	if !ok {
		return 0
	}

	return size
}

func getValueCost(key uint64) int16 {
	cost, ok := get(valueCosts, key)
	if !ok {
		return 0
	}

	return cost
}

func GetStaleKeys() []uint64 {
	now := time.Now().Unix()
	stale := make([]uint64, 0)

	valueTTLs.Mux.Lock()
	defer valueTTLs.Mux.Unlock()

	for key, valueTTL := range valueTTLs.M {
		if valueTTL < now {
			stale = append(stale, key)
		}
	}

	return stale
}

func getShardNum(key uint64) int {
	shardNum, ok := get(shardNums, key)
	if !ok {
		return 0
	}

	return shardNum
}

func SetPolicyAttr(key uint64, attr int64) {
	set(policyAttrs, key, attr)
}

func GetPolicyAttr(key uint64) (int64, bool) {
	attr, ok := get(policyAttrs, key)
	if !ok {
		return 0, false
	}

	return attr, true
}

func EnoughSpaceInShard(shardNum, size int) bool {
	volume := GetShardVolume(shardNum)
	return volume >= size
}

func GetShardVolume(shardNum int) int {
	volume, _ := get(shardVolumes, shardNum)
	return volume
}

func GetShardLock(shardNum int) *sync.RWMutex {
	return shardLocks[shardNum]
}

func GetKeyByMinPolicyAttr() uint64 {
	var hash uint64
	var currentCost int16

	minCost := int16(math.MaxInt16)
	minValue := int64(math.MaxInt64)

	policyAttrs.Mux.RLock()
	defer policyAttrs.Mux.RUnlock()

	for key, attr := range policyAttrs.M {
		currentCost = getValueCost(key)

		if attr <= minValue && currentCost <= minCost {
			minValue = attr
			minCost = currentCost
			hash = key
		}
	}

	return hash
}

func GetKeyByMaxPolicyAttr() uint64 {
	var hash uint64
	var maxCost int16
	var maxValue int64

	policyAttrs.Mux.RLock()
	defer policyAttrs.Mux.RUnlock()

	for key, attr := range policyAttrs.M {
		currentCost := getValueCost(key)

		if attr >= maxValue && currentCost >= maxCost {
			maxValue = attr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}

func GetKeyByMinIndex() uint64 {
	var hash uint64
	maxCost := int16(math.MaxInt16)
	minIndex := int(math.MaxInt64)

	policyAttrs.Mux.RLock()
	defer policyAttrs.Mux.RUnlock()

	for key, attr := range indexes.M {
		currentCost := getValueCost(key)

		if attr <= minIndex && currentCost <= maxCost {
			minIndex = attr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}
