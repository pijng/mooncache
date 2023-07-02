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

type Algorithm func() Variant

type PolicyService struct {
	getKeyByPolicy func() uint64
	updater        func(*keymaps.Keymaps, uint64)
	Variant        Variant
}

func Build(km *keymaps.Keymaps, variant Variant) *PolicyService {
	switch variant {
	case _LRU:
		return &PolicyService{
			getKeyByPolicy: km.GetKeyByMinPolicyAttr,
			updater:        updateKeyAttrByTime,
			Variant:        variant,
		}
	case _LFU:
		return &PolicyService{
			getKeyByPolicy: km.GetKeyByMinPolicyAttr,
			updater:        updateKeyAttrByCount,
			Variant:        variant,
		}
	case _MRU:
		return &PolicyService{
			getKeyByPolicy: km.GetKeyByMaxPolicyAttr,
			updater:        updateKeyAttrByTime,
			Variant:        variant,
		}
	case _MFU:
		return &PolicyService{
			getKeyByPolicy: km.GetKeyByMaxPolicyAttr,
			updater:        updateKeyAttrByCount,
			Variant:        variant,
		}
	case _FIFO:
		return &PolicyService{
			getKeyByPolicy: km.GetKeyByMinIndex,
			updater:        func(km *keymaps.Keymaps, u uint64) {},
			Variant:        variant,
		}
	default:
		return &PolicyService{}
	}
}

func (ps *PolicyService) UpdateKeyAttrByPolicy(km *keymaps.Keymaps, key uint64) {
	ps.updater(km, key)
}

func updateKeyAttrByTime(km *keymaps.Keymaps, key uint64) {
	_, ok := km.GetPolicyAttr(key)
	if !ok {
		return
	}

	currentTime := time.Now().Unix()
	km.SetPolicyAttr(key, currentTime)
}

func updateKeyAttrByCount(km *keymaps.Keymaps, key uint64) {
	count, ok := km.GetPolicyAttr(key)
	if !ok {
		return
	}

	km.SetPolicyAttr(key, count+1)
}

func (ps *PolicyService) EvictUntilCanFit(km *keymaps.Keymaps, size, shardNum int, del func(*keymaps.Keymaps, uint64)) {
	if km.EnoughSpaceInShard(shardNum, size) {
		return
	}

	hashedKey := ps.getKeyByPolicy()
	del(km, hashedKey)

	ps.EvictUntilCanFit(km, size, shardNum, del)
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
