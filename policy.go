package mooncache

import (
	"github.com/pijng/mooncache/internal/policy"
)

func LRU() policy.Policy {
	return policy.LRU
}

func MRU() policy.Policy {
	return policy.MRU
}

func LFU() policy.Policy {
	return policy.LFU
}

func MFU() policy.Policy {
	return policy.MFU
}

func FIFO() policy.Policy {
	return policy.FIFO
}
