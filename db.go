package gomemcache

import (
	"container/heap"
	"errors"
	"sync"
	"time"

	"github.com/Atoo35/go-mem-cache/src/constants"
	types "github.com/Atoo35/go-mem-cache/src/types"
)

type DB[T any] struct {
	mu                  sync.RWMutex
	capacity            int64
	data                map[string]types.Item[T]
	key_eviction_policy types.KeyEvictionPolicy
	pq                  priorityQueue[T]
}

type priorityQueue[T any] []types.Item[T]

func (pq priorityQueue[T]) Len() int {
	return len(pq)
}

func (pq *priorityQueue[T]) Less(i, j int) bool {
	return (*pq)[i].Priority < (*pq)[j].Priority
}

func (pq *priorityQueue[T]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}

func (pq *priorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*types.Item[T])
	item.Index = n
	*pq = append(*pq, *item)
}

func (pq *priorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Option[T any] func(*DB[T])

func New[T any](opts ...Option[T]) *DB[T] {
	db := &DB[T]{
		capacity:            constants.DefaultCapacity,
		data:                make(map[string]types.Item[T], constants.DefaultCapacity),
		key_eviction_policy: constants.DefaultKeyEvictionPolicy,
		pq:                  make(priorityQueue[T], 0),
	}
	for _, opt := range opts {
		opt(db)
	}
	return db
}

// Option functions
func WithKeyEvictionPolicy[T any](keyEvictionPolicy types.KeyEvictionPolicy) Option[T] {
	if !isValidEvictionPolicy(keyEvictionPolicy) {
		panic(constants.ErrInvalidKeyEvictionPolicy)
	}
	return func(db *DB[T]) {
		db.key_eviction_policy = keyEvictionPolicy
	}
}

func WithCapacity[T any](capacity int64) Option[T] {
	return func(db *DB[T]) {
		db.capacity = capacity
		db.data = make(map[string]types.Item[T], capacity)
	}
}

// Getter functions
func (db *DB[T]) GetCapacity() int64 {
	return db.capacity
}

func (db *DB[T]) GetKeyEvictionPolicy() types.KeyEvictionPolicy {
	return db.key_eviction_policy
}

// DB functions

func (c *DB[T]) GetKeys() []string {
	keys := make([]string, len(c.pq))
	for i, item := range c.pq {
		keys[i] = item.Key
	}
	return keys
}

func (db *DB[T]) Get(key string) (T, error) {
	db.mu.RLock()
	if item, ok := db.data[key]; ok {
		db.mu.RUnlock()
		// Check if the item has expired
		if item.Expiry <= time.Now().UnixNano() {
			// Remove the expired item
			heap.Remove(&db.pq, item.Index)
			delete(db.data, key)
			return item.Value, errors.New("item has expired")
		}

		// Update the last used time
		db.mu.Lock()
		item.Expiry = time.Now().Add(5 * time.Minute).UnixNano() // You can set a new expiration time here

		// Reorder the priority queue
		heap.Fix(&db.pq, item.Index)
		db.mu.Unlock()
		return item.Value, nil
	}

	return db.data[key].Value, errors.New("key not found")
}

func (db *DB[T]) SetWithExpiry(key string, value T, priority, expiry int64) {
	db.set(key, value, priority, expiry)
}

func (db *DB[T]) Set(key string, value T, priority int64) {
	db.set(key, value, priority, time.Now().Add(5*time.Minute).UnixNano())
}

func (db *DB[T]) Size() int64 {
	return db.size()
}

// private functions
func (db *DB[T]) set(key string, value T, priority, expiry int64) {
	db.mu.Lock()

	// Check if the cache exceeds its capacity
	if len(db.data) >= int(db.capacity) {
		// Find the item with the lowest priority
		// fmt.Println("setting key:", key, "value:", value, "priority:", priority, "expiry:", expiry)
		lowestPriorityItem := db.pq[0]
		lowestPriorityItems := make([]types.Item[T], 0)
		for _, item := range db.pq {
			if item.Priority == lowestPriorityItem.Priority {
				lowestPriorityItems = append(lowestPriorityItems, item) // Assign the address of item to lruItem
			}
		}
		// fmt.Println("lowerstItems", lowestPriorityItems)
		// Find the least recently used item
		lruItem := lowestPriorityItems[0]
		if len(lowestPriorityItems) > 1 {
			for _, item := range lowestPriorityItems {
				if item.Expiry < lruItem.Expiry {
					lruItem = item
				}
			}
		}

		// fmt.Println("lruItem:", lruItem.Key, "lruItem.Expiry:", lruItem.Expiry, "lruItem.Priority:", lruItem.Priority, "lruItem.Value:", lruItem.Value)

		// Remove the least recently used item
		heap.Remove(&db.pq, lruItem.Index)
		delete(db.data, lruItem.Key)
	}

	// Create a new CacheItem
	newItem := &types.Item[T]{
		Key:      key,
		Value:    value,
		Priority: priority,
		Expiry:   expiry,
	}

	// Add the new item to the priority queue and the item map
	heap.Push(&db.pq, newItem)
	db.data[key] = *newItem
	db.mu.Unlock()
}

func (db *DB[T]) size() int64 {
	db.mu.RLock()
	size := len(db.data)
	db.mu.RUnlock()
	return int64(size)
}

func isValidEvictionPolicy(evictionPolicy types.KeyEvictionPolicy) bool {
	switch evictionPolicy {
	case constants.LRUEviction:
		return true
	}
	return false
}
