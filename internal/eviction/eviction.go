package eviction

import (
	"sync"
	"time"

	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/shards"
)

var once sync.Once

func Build() {
	once.Do(func() {
		go worker()
	})
}

func worker() {
	timer := time.NewTicker(1 * time.Second)

	for {
		<-timer.C
		now := time.Now().Unix()
		evictOnTTL(now)
	}
}

func evictOnTTL(now int64) {
	staleKeys := keymaps.GetStaleKeys()

	for _, key := range staleKeys {
		shards.DelByHash(key)
	}
}
