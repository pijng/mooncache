package mooncache

import (
	"github.com/pijng/mooncache/internal/policy"
)

type Config struct {
	Algorithm    policy.Algorithm
	ShardSize    int
	ShardsAmount int8
}
