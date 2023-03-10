package policy

import (
	"time"

	"github.com/pijng/mooncache/internal/keymaps"
)

type Variant string

const (
	_LRU  = "LRU"
	_MRU  = "MRU"
	_LFU  = "LFU"
	_MFU  = "MFU"
	_FIFO = "FIFO"
)

type Policy func() Variant

var getKeyByPolicy func() uint64
var updater func(uint64)

func Build(variant Variant) {
	switch variant {
	case _LRU:
		getKeyByPolicy = keymaps.KeyByMinPolicyAttr
		updater = updateKeyAttrByTime
	case _LFU:
		getKeyByPolicy = keymaps.KeyByMinPolicyAttr
		updater = updateKeyAttrByCount
	case _MRU:
		getKeyByPolicy = keymaps.KeyByMaxPolicyAttr
		updater = updateKeyAttrByTime
	case _MFU:
		getKeyByPolicy = keymaps.KeyByMaxPolicyAttr
		updater = updateKeyAttrByCount
	case _FIFO:
		getKeyByPolicy = keymaps.KeyByMinIndex
		updater = func(u uint64) {}
	default:
		updater = func(u uint64) {}
	}
}

func UpdateKeyAttrByPolicy(key uint64) {
	updater(key)
}

func updateKeyAttrByTime(key uint64) {
	_, ok := keymaps.KeyPolicyAttr(key)
	if !ok {
		return
	}

	currentTime := time.Now().Unix()
	keymaps.SetKeyPolicyAttr(key, currentTime)
}

func updateKeyAttrByCount(key uint64) {
	count, ok := keymaps.KeyPolicyAttr(key)
	if !ok {
		return
	}

	keymaps.SetKeyPolicyAttr(key, count+1)
}

func EvictUntilCanFit(size, shardNum int, del func(uint64)) {
	if keymaps.EnoughSpaceInShard(shardNum, size) {
		return
	}

	hashedKey := getKeyByPolicy()
	del(hashedKey)

	EvictUntilCanFit(size, shardNum, del)
}

func LRU() Variant {
	return _LRU
}

func MRU() Variant {
	return _MRU
}

func LFU() Variant {
	return _LFU
}

func MFU() Variant {
	return _MFU
}

func FIFO() Variant {
	return _FIFO
}
