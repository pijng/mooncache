package shards

import (
	"math"
	"testing"

	"github.com/pijng/mooncache/internal/config"
	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/lib"
	"github.com/pijng/mooncache/internal/policy"
	"github.com/pijng/mooncache/internal/queue"
	"github.com/stretchr/testify/assert"
)

func buildCache() {
	configuration := config.Build(1<<10, 4, policy.MFU)
	keymaps.Build(configuration.ShardsAmount, configuration.ShardSize)
	Build(configuration.ShardsAmount)
	queue.Build()
	policy.Build(configuration.Policy())
}

func TestBuild(t *testing.T) {
	type args struct {
		amount int
	}
	tests := []struct {
		name string
		args args
	}{
		{"should build correct shards amount", args{amount: 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shards := Build(tt.args.amount)
			assert.Equal(t, tt.args.amount, len(*shards))
		})
	}
}

func TestSet(t *testing.T) {
	buildCache()

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int
		ttl   int64
	}
	tests := []struct {
		name string
		args args
	}{
		{"should set item to cache without errors", args{key: "item1", value: value{deep: "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"}, cost: 0, ttl: math.MaxInt64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Set(tt.args.key, tt.args.value, tt.args.cost, tt.args.ttl)
			assert.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	buildCache()

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int
		ttl   int64
	}
	wantValue := "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"item should be present in cache after set", args{key: "item1", value: value{deep: wantValue}, cost: 0, ttl: math.MaxInt64}, wantValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = Set(tt.args.key, tt.args.value, tt.args.cost, tt.args.ttl)
			got, err := Get(tt.args.key)
			assert.NoError(t, err)
			assert.Equal(t, wantValue, got.(value).deep)
		})
	}
}

func TestDel(t *testing.T) {
	buildCache()

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int
		ttl   int64
	}
	wantValue := "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"
	tests := []struct {
		name string
		args args
		err  error
	}{
		{"item should not be present in cache after delete", args{key: "item1337", value: value{deep: wantValue}, cost: 0, ttl: math.MaxInt64}, lib.ValueNotPresent()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = Set(tt.args.key, tt.args.value, tt.args.cost, tt.args.ttl)
			Del(tt.args.key)
			got, err := Get(tt.args.key)
			assert.Equal(t, tt.err, err)
			assert.Nil(t, got)
		})
	}
}
