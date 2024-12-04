package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	ID   int
	Name string
	Age  int
}

func TestCacheInstance_Set(t *testing.T) {
	cache := NewCacheInstance()

	tests := []struct {
		name  string
		key   string
		value interface{}
		ttl   time.Duration
	}{
		{
			name:  "string value",
			key:   "test-string",
			value: "hello world",
			ttl:   time.Second * 5,
		},
		{
			name:  "integer value",
			key:   "test-int",
			value: 42,
			ttl:   time.Second * 5,
		},
		{
			name: "struct value",
			key:  "test-struct",
			value: TestUser{
				ID:   1,
				Name: "John Doe",
				Age:  30,
			},
			ttl: time.Second * 5,
		},
		{
			name:  "nil value",
			key:   "test-nil",
			value: nil,
			ttl:   time.Second * 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Set(tt.key, tt.value, tt.ttl)

			value, exists := cache.Get(tt.key)
			require.True(t, exists)
			assert.Equal(t, tt.value, value)
		})
	}
}

func TestCacheInstance_Get(t *testing.T) {
	cache := NewCacheInstance()

	t.Run("get existing value", func(t *testing.T) {
		cache.Set("key1", "value1", time.Second*5)
		value, exists := cache.Get("key1")
		require.True(t, exists)
		assert.Equal(t, "value1", value)
	})

	t.Run("get non-existing value", func(t *testing.T) {
		value, exists := cache.Get("non-existing")
		require.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("get expired value", func(t *testing.T) {
		cache.Set("expired", "value", time.Millisecond*100)
		time.Sleep(time.Millisecond * 200)
		value, exists := cache.Get("expired")
		require.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("get struct and cast", func(t *testing.T) {
		user := TestUser{ID: 1, Name: "John", Age: 30}
		cache.Set("user", user, time.Second*5)

		value, exists := cache.Get("user")
		require.True(t, exists)

		castedUser, ok := value.(TestUser)
		require.True(t, ok)
		assert.Equal(t, user, castedUser)
	})
}

func TestCacheInstance_Del(t *testing.T) {
	cache := NewCacheInstance()

	t.Run("delete existing value", func(t *testing.T) {
		cache.Set("key1", "value1", time.Second*5)
		cache.Del("key1")
		_, exists := cache.Get("key1")
		assert.False(t, exists)
	})

	t.Run("delete non-existing value", func(t *testing.T) {
		cache.Del("non-existing")
		_, exists := cache.Get("non-existing")
		assert.False(t, exists)
	})
}

func TestCacheInstance_Pop(t *testing.T) {
	cache := NewCacheInstance()

	t.Run("pop existing value", func(t *testing.T) {
		cache.Set("key1", "value1", time.Second*5)
		value, exists := cache.Pop("key1")
		require.True(t, exists)
		assert.Equal(t, "value1", value)

		_, exists = cache.Get("key1")
		assert.False(t, exists)
	})

	t.Run("pop non-existing value", func(t *testing.T) {
		value, exists := cache.Pop("non-existing")
		require.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("pop expired value", func(t *testing.T) {
		cache.Set("expired", "value", time.Millisecond*100)
		time.Sleep(time.Millisecond * 200)
		value, exists := cache.Pop("expired")
		require.False(t, exists)
		assert.Nil(t, value)
	})
}

func TestCacheInstance_Len(t *testing.T) {
	cache := NewCacheInstance()

	assert.Equal(t, 0, cache.Len())

	cache.Set("key1", "value1", time.Second*5)
	assert.Equal(t, 1, cache.Len())

	cache.Set("key2", "value2", time.Second*5)
	assert.Equal(t, 2, cache.Len())

	cache.Del("key1")
	assert.Equal(t, 1, cache.Len())

	cache.Clear()
	assert.Equal(t, 0, cache.Len())
}

func TestCacheInstance_Clear(t *testing.T) {
	cache := NewCacheInstance()

	cache.Set("key1", "value1", time.Second*5)
	cache.Set("key2", "value2", time.Second*5)
	cache.Set("key3", "value3", time.Second*5)

	assert.Equal(t, 3, cache.Len())

	cache.Clear()

	assert.Equal(t, 0, cache.Len())

	_, exists := cache.Get("key1")
	assert.False(t, exists)
}

func TestCacheInstance_ConcurrentAccess(t *testing.T) {
	cache := NewCacheInstance()
	const goroutines = 100
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				key := strconv.Itoa(id*operationsPerGoroutine + j)

				cache.Set(key, id, time.Second*5)
				cache.Get(key)
				cache.Del(key)
				cache.Pop(key)
				cache.Len()
			}
		}(i)
	}

	wg.Wait()
}
