package hasher_test

import (
	"testing"

	"github.com/pijng/mooncache/internal/hasher"
)

func TestSumWithNum(t *testing.T) {
	type args struct {
		key          string
		shardsAmount int
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 int
	}{
		{"key should be 2546886805339723447 and num should be 2", args{key: "reports/sales", shardsAmount: 4}, 2546886805339723447, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := hasher.SumWithNum(tt.args.key, tt.args.shardsAmount)
			if got != tt.want {
				t.Errorf("SumWithNum() key = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SumWithNum() num = %v, want %v", got1, tt.want1)
			}
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
		{"key should be 565", args{key: "Pelinal Whitestrake is a cyborg proofs"}, 14388442895772204505},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasher.Sum(tt.args.key); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
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
			if got := hasher.JCH(tt.args.hashedKey, tt.args.shardsAmount); got != tt.want {
				t.Errorf("JCH() = %v, want %v", got, tt.want)
			}
		})
	}
}
