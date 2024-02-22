package types

type KeyEvictionPolicy string

type Item[T any] struct {
	Key      string
	Value    T
	Priority int64
	Expiry   int64
	Index    int // Index in the priority queue
}
