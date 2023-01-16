package queue_test

import (
	"sync"
	"testing"

	"github.com/pijng/mooncache/internal/queue"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	newQueue := queue.Build()

	tests := []struct {
		name string
	}{
		{"should return non nil queue"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEqual(t, nil, newQueue)
		})
	}
}

func TestSet(t *testing.T) {
	_ = queue.Build()

	type args struct {
		key uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{"should have transaction after set", args{key: 1337}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue.Set(tt.args.key)
			transaction := queue.Get(tt.args.key)
			assert.NotEqual(t, nil, transaction)
		})
	}
}

func TestSetDuplicate(t *testing.T) {
	_ = queue.Build()

	type args struct {
		key uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{"should not override transaction on duplicate set", args{key: 1337}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue.Set(tt.args.key)
			transaction1 := queue.Get(tt.args.key)

			queue.Set(tt.args.key)
			transaction2 := queue.Get(tt.args.key)

			assert.Equal(t, transaction1, transaction2)
		})
	}
}

func TestGet(t *testing.T) {
	_ = queue.Build()

	type args struct {
		key uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{"should have transaction after set", args{key: 1337}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue.Set(tt.args.key)
			transaction := queue.Get(tt.args.key)
			assert.NotEqual(t, nil, transaction)
		})
	}
}

func TestRelease(t *testing.T) {
	_ = queue.Build()

	type args struct {
		key         uint64
		transaction *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{"should have empty transaction after release", args{key: 1337, transaction: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue.Set(tt.args.key)
			queue.Release(tt.args.key)
			transaction := queue.Get(tt.args.key)
			assert.Equal(t, tt.args.transaction, transaction)
		})
	}
}
