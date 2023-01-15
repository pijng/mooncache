package mooncache

import (
	"time"

	"github.com/pijng/mooncache/internal/policy"
)

type ItemOptions struct {
	Cost int
	TTL  time.Duration
}

type policyContainer struct{}

var Policy policyContainer

func (p policyContainer) LRU() policy.Policy {
	return policy.LRU
}

func (p policyContainer) MRU() policy.Policy {
	return policy.MRU
}

func (p policyContainer) LFU() policy.Policy {
	return policy.LFU
}

func (p policyContainer) MFU() policy.Policy {
	return policy.MFU
}

func (p policyContainer) FIFO() policy.Policy {
	return policy.FIFO
}
