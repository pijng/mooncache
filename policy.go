package mooncache

import (
	"github.com/pijng/mooncache/internal/policy"
)

func LRU() policy.Algorithm {
	return policy.LRU
}

func MRU() policy.Algorithm {
	return policy.MRU
}

func LFU() policy.Algorithm {
	return policy.LFU
}

func MFU() policy.Algorithm {
	return policy.MFU
}

func FIFO() policy.Algorithm {
	return policy.FIFO
}
