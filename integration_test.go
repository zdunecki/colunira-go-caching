package main

import (
	"gosession/caching/cache"
	"gosession/caching/database"
	"sync"
	"testing"
	"time"
)

func TestSingleAccessCaching(t *testing.T) {
	/// Given

	id := "test-id"
	data := "example data"
	d := database.Database{
		Data: map[string]interface{}{
			id: data,
		},
	}
	expiresAfter, _ := time.ParseDuration("5s")
	c := cache.New(expiresAfter)

	/// When

	_, cached := c.Get(id, d.GetById)

	t.Run("should not find data in cache before it was cached", func(t *testing.T) {
		if cached {
			t.Fatal("Data should not be available in cache before using c.Set")
		}
	})

	expectedResult := d.GetById(id)
	cachedResult, cached := c.Get(id, d.GetById)

	t.Run("should find expected data in cache after it was cached", func(t *testing.T) {
		if !cached {
			t.Fatal("Data should be available in cache after using c.Set")
		}
		if expectedResult != cachedResult {
			t.Fatalf("Cached and database result should be the same. Expected: %v, got: %v", expectedResult, cachedResult)
		}
	})

	time.Sleep(expiresAfter)
	_, expiredFound := c.Get(id, d.GetById)

	t.Run("should not find data in cache after it expired", func(t *testing.T) {
		if expiredFound {
			t.Fatalf("Data should not be available in the cache after expiration time passed")
		}
	})
}

type Result struct {
	data  interface{}
	cached bool
}

func TestMultiAccessCaching(t *testing.T) {
	/// Given

	id := "test-id"
	data := "example data"
	d := database.Database{
		Data: map[string]interface{}{
			id: data,
		},
	}
	expiresAfter, _ := time.ParseDuration("15s")
	c := cache.New(expiresAfter)

	numWorkers := 100

	channel := make(chan Result, numWorkers)
	results := make([]Result, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	/// When

	for range results {
		go func() {
			result, found := c.Get(id, d.GetById)

			channel <- Result{result, found}

			wg.Done()
		}()
	}

	wg.Wait()

	for i := range results {
		results[i] = <-channel
	}

	// Then

	notCachedCount := 0
	for _, result := range results {
		if !result.cached {
			notCachedCount++
		}
	}

	t.Run("should cache data after one request to database", func(t* testing.T) {
		if notCachedCount != 1 {
			t.Fatalf("Data in cache should be unavailable only once, actual: %v", notCachedCount)
		}
	})

}
