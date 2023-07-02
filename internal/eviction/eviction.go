package eviction

import (
	"time"

	"github.com/pijng/mooncache/internal/keymaps"
	"github.com/pijng/mooncache/internal/shards"
)

func Run(shards shards.Shards, km *keymaps.Keymaps) {
	go worker(shards, km)
}

func worker(shards shards.Shards, km *keymaps.Keymaps) {
	timer := time.NewTicker(1 * time.Second)

	for {
		<-timer.C
		now := time.Now().Unix()
		evictOnTTL(shards, km, now)
	}
}

func evictOnTTL(shards shards.Shards, km *keymaps.Keymaps, now int64) {
	staleKeys := km.GetStaleKeys()

	for _, key := range staleKeys {
		shards.DelByHash(km, key)
	}
}
