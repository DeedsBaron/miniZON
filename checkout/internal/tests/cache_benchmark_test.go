package tests

import (
	"context"
	"fmt"
	"testing"

	cache2 "route256/libs/cache"
	"route256/libs/logger"
)

func BenchmarkLRUCache(b *testing.B) {
	logger.Init(true, "DEBUG")
	ctx := context.Background()
	// Create a new LRU cache with a capacity of 100 and 4 buckets
	cache := cache2.NewCache[string](ctx, 0, 10, 1000)

	// Add some key-value pairs to the cache
	cache.Set(ctx, "key1", "value1")
	cache.Set(ctx, "key2", "value2")
	cache.Set(ctx, "key3", "value3")
	cache.Set(ctx, "key4", "value4")
	cache.Set(ctx, "key5", "value5")

	// Benchmark the cache by performing 1000 lookups
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, "key1")
		cache.Get(ctx, "key2")
		cache.Get(ctx, "key3")
		cache.Get(ctx, "key4")
		cache.Get(ctx, "key5")
	}
}

func TestLRUCache(t *testing.T) {
	// Run the benchmark and print the results
	result := testing.Benchmark(BenchmarkLRUCache)
	fmt.Println(result.String())
}
