package hasher_test

import (
	"testing"

	"github.com/pijng/mooncache/internal/hasher"
	"github.com/stretchr/testify/assert"
)

func TestSumWithNum(t *testing.T) {
	type args struct {
		key          string
		shardsAmount int
	}
	tests := []struct {
		name    string
		args    args
		wantSum uint64
		wantNum int
	}{
		{"key should be 2546886805339723447 and num should be 2", args{key: "reports/sales", shardsAmount: 4}, 2546886805339723447, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, num := hasher.SumWithNum(tt.args.key, tt.args.shardsAmount)
			assert.Equal(t, tt.wantSum, sum)
			assert.Equal(t, tt.wantNum, num)
		})
	}
}

func TestSum(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"key should be 14388442895772204505", args{key: "Pelinal Whitestrake is a cyborg proofs"}, 14388442895772204505},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := hasher.Sum(tt.args.key)
			assert.Equal(t, tt.want, sum)
		})
	}
}

func TestJCH(t *testing.T) {
	type args struct {
		hashedKey    uint64
		shardsAmount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"num should be 6", args{hashedKey: 5577006791947779411, shardsAmount: 8}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num := hasher.JCH(tt.args.hashedKey, tt.args.shardsAmount)
			assert.Equal(t, tt.want, num)
		})
	}
}
