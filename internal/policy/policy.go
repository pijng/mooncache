package policy

import (
	"sync"
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
var UpdateKeyAttrByPolicy func(uint64)
var once sync.Once

func Build(variant Variant) {
	once.Do(func() {
		switch variant {
		case _LRU:
			getKeyByPolicy = keymaps.GetKeyByMinPolicyAttr
			UpdateKeyAttrByPolicy = updateKeyAttrByTime
		case _LFU:
			getKeyByPolicy = keymaps.GetKeyByMinPolicyAttr
			UpdateKeyAttrByPolicy = updateKeyAttrByCount
		case _MRU:
			getKeyByPolicy = keymaps.GetKeyByMaxPolicyAttr
			UpdateKeyAttrByPolicy = updateKeyAttrByTime
		case _MFU:
			getKeyByPolicy = keymaps.GetKeyByMaxPolicyAttr
			UpdateKeyAttrByPolicy = updateKeyAttrByCount
		case _FIFO:
			getKeyByPolicy = keymaps.GetKeyByMinIndex
		default:
		}
	})
}

func updateKeyAttrByTime(key uint64) {
	_, ok := keymaps.GetKeyPolicyAttr(key)
	if !ok {
		return
	}

	currentTime := time.Now().Unix()
	keymaps.SetKeyPolicyAttr(key, currentTime)
}

func updateKeyAttrByCount(key uint64) {
	count, ok := keymaps.GetKeyPolicyAttr(key)
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
