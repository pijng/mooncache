package queue

import (
	"sync"
)

type Queue struct {
	hashMap map[uint64]*sync.WaitGroup
	mux     sync.RWMutex
}

var queue Queue

func Build() *Queue {
	queue = Queue{
		hashMap: make(map[uint64]*sync.WaitGroup),
	}

	return &queue
}

func Set(key uint64) {
	queue.mux.Lock()
	defer queue.mux.Unlock()

	_, ok := queue.hashMap[key]
	if ok {
		return
	}

	queue.hashMap[key] = &sync.WaitGroup{}
	queue.hashMap[key].Add(1)
}

func Get(key uint64) *sync.WaitGroup {
	queue.mux.RLock()
	defer queue.mux.RUnlock()

	return queue.hashMap[key]
}

func Release(key uint64) {
	queue.mux.Lock()
	defer queue.mux.Unlock()

	if transaction, ok := queue.hashMap[key]; ok {
		transaction.Done()
		delete(queue.hashMap, key)
	}
}
