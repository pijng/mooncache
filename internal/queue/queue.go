package queue

import (
	"sync"
)

type Queue struct {
	hashMap map[uint64]*sync.WaitGroup
	mux     sync.RWMutex
}

func Build() *Queue {
	return &Queue{
		hashMap: make(map[uint64]*sync.WaitGroup),
	}
}

func (q *Queue) Set(key uint64) {
	q.mux.Lock()
	defer q.mux.Unlock()

	_, ok := q.hashMap[key]
	if ok {
		return
	}

	q.hashMap[key] = &sync.WaitGroup{}
	q.hashMap[key].Add(1)
}

func (q *Queue) Get(key uint64) *sync.WaitGroup {
	q.mux.RLock()
	defer q.mux.RUnlock()

	return q.hashMap[key]
}

func (q *Queue) Release(key uint64) {
	q.mux.Lock()
	defer q.mux.Unlock()

	if transaction, ok := q.hashMap[key]; ok {
		transaction.Done()
		delete(q.hashMap, key)
	}
}
