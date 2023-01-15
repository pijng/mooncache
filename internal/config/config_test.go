package config_test

import (
	"reflect"
	"testing"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
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

func TestConfig(t *testing.T) {
	config.Build(1<<10, 4, policy.FIFO)

	tests := []struct {
		name string
		want *config.Configuration
	}{
		{"should return valid config", &config.Configuration{ShardSize: 1 << 10, ShardsAmount: 4, Policy: policy.FIFO}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration := config.Config()

			assert.Equal(t, tt.want.ShardSize, configuration.ShardSize)
			assert.Equal(t, tt.want.ShardsAmount, configuration.ShardsAmount)
			assert.Equal(t, tt.want.Policy(), configuration.Policy())
		})
	}
}

func TestShardSize(t *testing.T) {
	config.Build(1<<10, 0, nil)

	tests := []struct {
		name string
		want int
	}{
		{"should be 1024", 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.ShardSize(); got != tt.want {
				t.Errorf("ShardSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShardsAmount(t *testing.T) {
	config.Build(0, 8, nil)

	tests := []struct {
		name string
		want int
	}{
		{"should be 8", 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.ShardsAmount(); got != tt.want {
				t.Errorf("ShardsAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicy(t *testing.T) {
	config.Build(0, 0, policy.MFU)

	tests := []struct {
		name string
		want policy.Policy
	}{
		{"should be MFU", policy.MFU},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.Policy(); !reflect.DeepEqual(got(), tt.want()) {
				t.Errorf("Policy() = %v, want %v", got, tt.want)
			}
		})
	}
}
