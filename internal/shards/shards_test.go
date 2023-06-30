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

func buildCache(policyVariant policy.Policy) {
	configuration := config.Build(1<<8, 1, policyVariant)
	keymaps.Build(configuration.ShardsAmount, configuration.ShardSize)
	Build(configuration.ShardsAmount)
	queue.Build()
	if policyVariant != nil {
		policy.Build(configuration.Policy())
	}
}

func TestBuild(t *testing.T) {
	type args struct {
		amount int8
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

			assert.Equal(t, tt.args.amount, int8(len(*shards)))
		})
	}
}

func TestSet(t *testing.T) {
	buildCache(policy.MFU)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int16
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

func TestSetToFullShardWithPolicy(t *testing.T) {
	buildCache(policy.MFU)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int16
		ttl   int64
	}
	tests := []struct {
		name  string
		item1 args
		item2 args
	}{
		{"should set item to cache without errors when policy is set and item size is fittable",
			args{
				key:   "item1",
				value: value{deep: "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"},
				cost:  0, ttl: math.MaxInt64,
			},
			args{
				key:   "item2",
				value: value{deep: "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"},
				cost:  0, ttl: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err1 := Set(tt.item1.key, tt.item1.value, tt.item1.cost, tt.item1.ttl)
			err2 := Set(tt.item2.key, tt.item2.value, tt.item2.cost, tt.item2.ttl)

			assert.NoError(t, err1)
			assert.NoError(t, err2)
		})
	}
}

func TestSetToFullShardWithoutPolicy(t *testing.T) {
	buildCache(nil)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int16
		ttl   int64
	}
	tests := []struct {
		name  string
		item1 args
		item2 args
	}{
		{"should set item to cache without errors when policy is set and item size is fittable",
			args{
				key:   "item1",
				value: value{deep: "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"},
				cost:  0, ttl: math.MaxInt64,
			},
			args{
				key:   "item2",
				value: value{deep: "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"},
				cost:  0, ttl: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err1 := Set(tt.item1.key, tt.item1.value, tt.item1.cost, tt.item1.ttl)
			err2 := Set(tt.item2.key, tt.item2.value, tt.item2.cost, tt.item2.ttl)

			assert.NoError(t, err1)
			assert.Error(t, err2)
		})
	}
}

func TestGet(t *testing.T) {
	buildCache(policy.MFU)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value value
		cost  int16
		ttl   int64
	}
	wantValue := "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"

	tests := []struct {
		name      string
		args      args
		shouldSet bool
		want      interface{}
		err       error
	}{
		{"item should be present in cache after set", args{key: "item1", value: value{deep: wantValue}, cost: 0, ttl: math.MaxInt64}, true, wantValue, nil},
		{"item should not be present in cache without set", args{key: "item3", value: value{}, cost: 0, ttl: math.MaxInt64}, false, nil, lib.ValueNotPresent()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSet {
				_ = Set(tt.args.key, tt.args.value, tt.args.cost, tt.args.ttl)
			}

			got, err := Get(tt.args.key)
			assert.Equal(t, tt.err, err)

			if got != nil {
				assert.Equal(t, tt.args.value.deep, got.(value).deep)
			}
		})
	}
}

func TestDel(t *testing.T) {
	buildCache(policy.MFU)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int16
		ttl   int64
	}
	wantValue := "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"
	tests := []struct {
		name string
		args args
		err  error
	}{
		{"item should not be present in cache after delete by string key", args{key: "item1337", value: value{deep: wantValue}, cost: 0, ttl: math.MaxInt64}, lib.ValueNotPresent()},
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

func TestDelByHash(t *testing.T) {
	buildCache(policy.MFU)

	type value struct {
		deep string
	}
	type args struct {
		key   string
		value interface{}
		cost  int16
		ttl   int64
	}
	wantValue := "01001111 01101000 00111111 00100000 01011001 01101111 01110101 00100111 01110010 01100101 00100000 01000001 01110000 01110000 01110010 01101111 01100001 01100011 01101000 01101001 01101110 01100111 00100000 01001101 01100101 00111111"
	tests := []struct {
		name      string
		args      args
		shouldSet bool
		hashedKey uint64
		err       error
	}{
		{"item should not be present in cache after delete by hashed key", args{key: "Pelinal Whitestrake is a cyborg proofs", value: value{deep: wantValue}, cost: 0, ttl: math.MaxInt64}, true, 14388442895772204505, lib.ValueNotPresent()},
		{"return if there is no item with such hashed key in cache", args{key: "407 1505"}, false, 4342668836759772583, lib.ValueNotPresent()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSet {
				_ = Set(tt.args.key, tt.args.value, tt.args.cost, tt.args.ttl)
			}

			DelByHash(tt.hashedKey)
			got, err := Get(tt.args.key)

			assert.Equal(t, tt.err, err)
			assert.Nil(t, got)
		})
	}
}
