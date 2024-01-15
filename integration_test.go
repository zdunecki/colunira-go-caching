package main

import (
	"gosession/caching/cache"
	"gosession/caching/database"
	"testing"
	"time"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestSingleAccessCaching(t *testing.T) {
	/// Setup
	id := "test-id"
	data := "example data"
	d := database.Database{
		Data: map[string]interface{}{
			id: data,
		},
	}
	expires, _ := time.ParseDuration("5s")
	c := cache.New(expires)

	/// Test

	// First Get should not find a value
	_, notCachedFound := c.Get(id)
	if notCachedFound {
		t.Fatal("Data should not be available in cache before using c.Set")
	}

	expectedResult, _ := d.GetById(id)
	c.Set(id, expectedResult)

	// Next Get should find a value until expired
	cachedResult, cachedFound := c.Get(id)
	if !cachedFound {
		t.Fatal("Data should be available in cache after using c.Set")
	}

	// Cached and database result should be the same
	if expectedResult != cachedResult {
		t.Fatalf("Cached and database result should be the same. Expected: %v, got: %v", expectedResult, cachedResult)
	}

	time.Sleep(expires)
	_, expiredFound := c.Get(id)
	// After expiration time the value should be expired and therefore not found in the cache
	if expiredFound {
		t.Fatalf("Data should not be available in the cache after expiration time passed")
	}
}

type Result struct {
	data interface{}
	found bool
}

func getWithCaching(ch chan Result, c cache.CacheInterface, d database.Database, id string) {
	result, found := c.Get(id)

	if !found {
		result, _ = d.GetById(id)
		c.Set(id, result)
	}

	ch <- Result{result, found}
}

func TestMultiAccessCaching(t *testing.T) {
	/// Setup
	id := "test-id"
	data := "example data"
	d := database.Database{
		Data: map[string]interface{}{
			id: data,
		},
	}
	expires, _ := time.ParseDuration("5s")
	c := cache.New(expires)

	/// Test

	channels := make([]chan Result, 10)
	results := make([]Result, 10)

	for i := range channels {
		channels[i] = make(chan Result)
		go getWithCaching(channels[i], c, d, id)
	}

	for i := range results {
		results[i] = <- channels[i]
	}

	notFoundCount := 0
	for _, result := range results {
		if !result.found {
			notFoundCount++
		}
	}

	if notFoundCount != 1 {
		t.Fatalf("Data in cache should be unavailable only once, actual: %v", notFoundCount)
	}
}
