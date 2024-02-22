package constants

import types "github.com/Atoo35/go-mem-cache/src/types"

const (
	DefaultCapacity                                  = 10
	DefaultKeyEvictionPolicy types.KeyEvictionPolicy = "LRU"

	LRUEviction types.KeyEvictionPolicy = "LRU"

	ErrInvalidKeyEvictionPolicy = "invalid key eviction policy"
)
