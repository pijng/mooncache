package lib_test

import (
	"testing"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/stretchr/testify/assert"
)

func TestCantFitInShard(t *testing.T) {
	type args struct {
		shardSize int
		shardNum  int
		size      int
	}

	config := config.Build(1<<10, 1, nil)
	keymaps.Build(config.ShardsAmount, config.ShardSize)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"should be true when value size is bigger than shard size", args{shardSize: config.ShardSize, shardNum: 0, size: 1 << 11}, true},
		{"should be false when value size is less than shard size", args{shardSize: config.ShardSize, shardNum: 0, size: 1 << 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.CantFitInShard(tt.args.shardSize, tt.args.shardNum, tt.args.size); got != tt.want {
				t.Errorf("CantFitInShard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCantFitInShardWithPolicy(t *testing.T) {
	type args struct {
		shardSize int
		shardNum  int
		size      int
	}

	config := config.Build(1<<10, 1, policy.LRU)
	keymaps.Build(config.ShardsAmount, config.ShardSize)
	keymaps.AddKey(0, 0, 0, 1<<10, 0, 0)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"should be true when value size is bigger than shard size", args{shardSize: config.ShardSize, shardNum: 0, size: 1 << 11}, true},
		{"should be false when value size is less than shard size and policy present and shard doesn't have enough space", args{shardSize: config.ShardSize, shardNum: 0, size: 1 << 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.CantFitInShard(tt.args.shardSize, tt.args.shardNum, tt.args.size); got != tt.want {
				t.Errorf("CantFitInShard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueSize(t *testing.T) {
	type args struct {
		value interface{}
	}

	type value struct {
		id     int
		number string
		total  int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"value size should be 12", args{value: value{id: 1, number: "1-1", total: 1200}}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, 12, lib.ValueSize(tt.args.value))
		})
	}
}
