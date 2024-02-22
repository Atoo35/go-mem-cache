package gomemcache

import (
	"strconv"
	"testing"

	types "github.com/Atoo35/go-mem-cache/src/types"
	"github.com/stretchr/testify/assert"
)

func TestDB_New(t *testing.T) {
	db := New[int64]()

	assert.Equal(t, int64(10), db.GetCapacity())
}

func TestDB_NewWithCapacity(t *testing.T) {
	db := New[int64](WithCapacity[int64](20))

	assert.Equal(t, int64(20), db.GetCapacity())
}

func TestDB_NewWithKeyEvictionPolicy(t *testing.T) {
	db := New[int64](WithKeyEvictionPolicy[int64](types.KeyEvictionPolicy("LRU")))

	assert.Equal(t, types.KeyEvictionPolicy("LRU"), db.GetKeyEvictionPolicy())
}

func TestDB_NewWithKeyEvictionPolicyPanic(t *testing.T) {
	assert.Panics(t, func() {
		New[int64](WithKeyEvictionPolicy[int64](types.KeyEvictionPolicy("test")))
	})
}

func TestDB_Set(t *testing.T) {
	testCases := []struct {
		key         string
		value       string
		expectedRes string
		expectedErr error
	}{
		{
			key:         "key",
			value:       "value",
			expectedRes: "value",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		db := New[string]()

		db.Set(tc.key, tc.value, 1)
		res, err := db.Get(tc.key)
		assert.Equal(t, tc.expectedRes, res)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestDB_Set11Items(t *testing.T) {
	db := New[int64](WithCapacity[int64](10))

	for i := 0; i < 11; i++ {
		db.Set("key"+strconv.Itoa(i), int64(i), 1)
	}

	assert.Equal(t, int64(1), db.data["key1"].Value)
}

func BenchmarkSetWithExpiry(b *testing.B) {
	db := New[int64](
		WithCapacity[int64](1000), // Adjust capacity as needed
	)

	// Run the SetWithExpiry method b.N times
	for i := 0; i < b.N; i++ {
		db.SetWithExpiry("key", 123, 5, 1000) // Adjust key, value, priority, and expiry as needed
	}
}

// BenchmarkGet benchmarks the Get method.
func BenchmarkGet(b *testing.B) {
	db := New[int64](
		WithCapacity[int64](1000), // Adjust capacity as needed
	)
	db.Set("key", 123, 5) // Set a key-value pair to retrieve

	// Run the Get method b.N times
	for i := 0; i < b.N; i++ {
		_, _ = db.Get("key") // Adjust key as needed
	}
}
