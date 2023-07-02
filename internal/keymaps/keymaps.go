package keymaps

import (
	"math"
	"sync"
	"time"
)

type Keymaps struct {
	shardVolumes hashmap[int, int]
	indexes      hashmap[uint64, int]
	shardNums    hashmap[uint64, int]
	valueSizes   hashmap[uint64, int]
	valueCosts   hashmap[uint64, int16]
	valueTTLs    hashmap[uint64, int64]
	policyAttrs  hashmap[uint64, int64]
	shardLocks   map[int]*sync.RWMutex
}

type hashmap[K int | uint64, T int | int16 | int64] struct {
	Mux *sync.RWMutex
	M   map[K]T
}

func Build(shardsAmount int8, shardSize int) *Keymaps {
	keymaps := &Keymaps{
		shardVolumes: hashmap[int, int]{
			Mux: &sync.RWMutex{},
			M:   make(map[int]int),
		},
		indexes: hashmap[uint64, int]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int),
		},
		shardNums: hashmap[uint64, int]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int),
		},
		valueSizes: hashmap[uint64, int]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int),
		},
		valueCosts: hashmap[uint64, int16]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int16),
		},
		valueTTLs: hashmap[uint64, int64]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int64),
		},
		policyAttrs: hashmap[uint64, int64]{
			Mux: &sync.RWMutex{},
			M:   make(map[uint64]int64),
		},
		shardLocks: make(map[int]*sync.RWMutex),
	}

	for n := 0; n < int(shardsAmount); n++ {
		keymaps.shardVolumes.M[n] = shardSize
		keymaps.shardLocks[n] = &sync.RWMutex{}
	}

	return keymaps
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

func (km *Keymaps) AddKey(key uint64, index, shardNum, size int, cost int16, ttl int64) {
	set(km.indexes, key, index)
	set(km.valueCosts, key, cost)
	set(km.valueTTLs, key, ttl)

	km.decrementShardVolume(shardNum, size)
	set(km.valueSizes, key, size)
	set(km.shardNums, key, shardNum)
	set(km.policyAttrs, key, 0)
}

func (km *Keymaps) DelKey(key uint64) {
	remove(km.indexes, key)
	remove(km.valueCosts, key)
	remove(km.valueTTLs, key)

	km.incrementShardVolume(key)

	remove(km.valueSizes, key)
	remove(km.shardNums, key)
	remove(km.policyAttrs, key)
}

func (km *Keymaps) decrementShardVolume(shardNum, size int) {
	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(km.shardVolumes, shardNum)
	set(km.shardVolumes, shardNum, currentVolume-size)
}

func (km *Keymaps) incrementShardVolume(key uint64) {
	shardNum, ok := get(km.shardNums, key)
	if !ok {
		return
	}

	size, ok := get(km.valueSizes, key)
	if !ok {
		return
	}

	// shardVolume at the given shardNum is always present
	currentVolume, _ := get(km.shardVolumes, shardNum)
	set(km.shardVolumes, shardNum, currentVolume+size)
}

func (km *Keymaps) GetIndex(key uint64) (int, bool) {
	index, ok := get(km.indexes, key)
	if !ok {
		return 0, false
	}

	return index, true
}

func getValueCost(valueCosts hashmap[uint64, int16], key uint64) int16 {
	cost, ok := get(valueCosts, key)
	if !ok {
		return 0
	}

	return cost
}

func (km *Keymaps) GetStaleKeys() []uint64 {
	now := time.Now().Unix()
	stale := make([]uint64, 0)

	km.valueTTLs.Mux.Lock()
	defer km.valueTTLs.Mux.Unlock()

	for key, valueTTL := range km.valueTTLs.M {
		if valueTTL < now {
			stale = append(stale, key)
		}
	}

	return stale
}

func (km *Keymaps) SetPolicyAttr(key uint64, attr int64) {
	set(km.policyAttrs, key, attr)
}

func (km *Keymaps) GetPolicyAttr(key uint64) (int64, bool) {
	attr, ok := get(km.policyAttrs, key)
	if !ok {
		return 0, false
	}

	return attr, true
}

func (km *Keymaps) EnoughSpaceInShard(shardNum, size int) bool {
	volume := km.GetShardVolume(shardNum)
	return volume >= size
}

func (km *Keymaps) GetShardVolume(shardNum int) int {
	volume, _ := get(km.shardVolumes, shardNum)
	return volume
}

func (km *Keymaps) GetShardLock(shardNum int) *sync.RWMutex {
	return km.shardLocks[shardNum]
}

func (km *Keymaps) GetKeyByMinPolicyAttr() uint64 {
	var hash uint64
	var currentCost int16

	minCost := int16(math.MaxInt16)
	minValue := int64(math.MaxInt64)

	km.policyAttrs.Mux.RLock()
	defer km.policyAttrs.Mux.RUnlock()

	for key, attr := range km.policyAttrs.M {
		currentCost = getValueCost(km.valueCosts, key)

		if attr <= minValue && currentCost <= minCost {
			minValue = attr
			minCost = currentCost
			hash = key
		}
	}

	return hash
}

func (km *Keymaps) GetKeyByMaxPolicyAttr() uint64 {
	var hash uint64
	var maxCost int16
	var maxValue int64

	km.policyAttrs.Mux.RLock()
	defer km.policyAttrs.Mux.RUnlock()

	for key, attr := range km.policyAttrs.M {
		currentCost := getValueCost(km.valueCosts, key)

		if attr >= maxValue && currentCost >= maxCost {
			maxValue = attr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}

func (km *Keymaps) GetKeyByMinIndex() uint64 {
	var hash uint64
	maxCost := int16(math.MaxInt16)
	minIndex := int(math.MaxInt64)

	km.policyAttrs.Mux.RLock()
	defer km.policyAttrs.Mux.RUnlock()

	for key, attr := range km.indexes.M {
		currentCost := getValueCost(km.valueCosts, key)

		if attr <= minIndex && currentCost <= maxCost {
			minIndex = attr
			maxCost = currentCost
			hash = key
		}
	}

	return hash
}
