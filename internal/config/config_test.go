package config_test

import (
	"testing"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/stretchr/testify/assert"
)

func TestAssignConfig(t *testing.T) {
	type args struct {
		shardSize    int
		shardsAmount int
		policy       policy.Policy
	}

	configArgs := args{
		shardSize:    1 << 8,
		shardsAmount: 4,
		policy:       policy.LFU,
	}

	tests := []struct {
		name string
		args args
	}{
		{"config fields should be set", configArgs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.Build(tt.args.shardSize, tt.args.shardsAmount, tt.args.policy)
			assert.Equal(t, tt.args.shardsAmount, config.ShardsAmount)
			assert.Equal(t, tt.args.shardSize, config.ShardSize)
			assert.Equal(t, tt.args.policy(), config.Policy())
		})
	}
}
