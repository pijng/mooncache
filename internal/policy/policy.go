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

type policyService struct {
	getKeyByPolicy func() uint64
	updater        func(uint64)
}

var service policyService

func Build(variant Variant) {
	switch variant {
	case _LRU:
		service = policyService{
			getKeyByPolicy: keymaps.GetKeyByMinPolicyAttr,
			updater:        updateKeyAttrByTime,
		}
	case _LFU:
		service = policyService{
			getKeyByPolicy: keymaps.GetKeyByMinPolicyAttr,
			updater:        updateKeyAttrByCount,
		}
	case _MRU:
		service = policyService{
			getKeyByPolicy: keymaps.GetKeyByMaxPolicyAttr,
			updater:        updateKeyAttrByTime,
		}
	case _MFU:
		service = policyService{
			getKeyByPolicy: keymaps.GetKeyByMaxPolicyAttr,
			updater:        updateKeyAttrByCount,
		}
	case _FIFO:
		service = policyService{
			getKeyByPolicy: keymaps.GetKeyByMinIndex,
			updater:        func(u uint64) {},
		}
	default:
	}
}

func UpdateKeyAttrByPolicy(key uint64) {
	service.updater(key)
}

func updateKeyAttrByTime(key uint64) {
	_, ok := keymaps.GetPolicyAttr(key)
	if !ok {
		return
	}

	currentTime := time.Now().Unix()
	keymaps.SetPolicyAttr(key, currentTime)
}

func updateKeyAttrByCount(key uint64) {
	count, ok := keymaps.GetPolicyAttr(key)
	if !ok {
		return
	}

	keymaps.SetPolicyAttr(key, count+1)
}

func EvictUntilCanFit(size, shardNum int, del func(uint64)) {
	if keymaps.EnoughSpaceInShard(shardNum, size) {
		return
	}

	hashedKey := service.getKeyByPolicy()
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
