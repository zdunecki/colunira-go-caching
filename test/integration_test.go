package cache_test

import (
	"gosession/caching/cache"
	"sync"
	"testing"
	"time"
)

func TestSingleAccessCaching(t *testing.T) {
	/// Given

	id := "test-id"
	data := "example data"
	d := MockDatabase{
		Data: map[string]interface{}{
			id: data,
		},
	}
	expiresAfter := time.Second
	c := cache.New(expiresAfter)

	/// When

	_, cached := c.Get(id, d.GetById)

	t.Run("should not find data in cache on first access", func(t *testing.T) {
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
	data   interface{}
	cached bool
}

func TestMultiAccessCaching(t *testing.T) {
	/// Given

	id := "test-id"
	expectedData := "example data"
	d := MockDatabase{
		Data: map[string]interface{}{
			id: expectedData,
		},
		OperationTime: 500 * time.Millisecond,
	}
	expiresAfter := 15 * time.Second
	c := cache.New(expiresAfter)

	numWorkers := 100

	channel := make(chan Result, numWorkers)
	results := make([]Result, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	/// When

	for range results {
		go func() {
			result, cached := c.Get(id, d.GetById)

			channel <- Result{result, cached}

			wg.Done()
		}()
	}

	wg.Wait()

	for i := range results {
		results[i] = <-channel
	}

	notCachedCount := 0
	validResult := true
	for _, result := range results {
		if !result.cached {
			notCachedCount++
		}
		if result.data != expectedData {
			validResult = false
		}
	}

	t.Run("should always get valid data from cache and database", func(t* testing.T) {
		if !validResult {
			t.Fatalf("Data from cache is not valid")
		}
	})
	t.Run("should cache data after one request to database", func(t *testing.T) {
		if notCachedCount != 1 {
			t.Fatalf("Data in cache should be unavailable only once, actual: %v", notCachedCount)
		}
	})

}

func TestMultiAccessCachingExpiration(t *testing.T) {
	/// Given

	id := "test-id"
	expectedData := "example data"
	d := MockDatabase{
		Data: map[string]interface{}{
			id: expectedData,
		},
		OperationTime: 500 * time.Millisecond,
	}
	expiresAfter := time.Second
	c := cache.New(expiresAfter)

	numWorkers := 100

	channel := make(chan Result, numWorkers)
	results := make([]Result, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	/// When

	for i := range results {
		if i == numWorkers/2 {
			time.Sleep(2 * time.Second)
		}
		go func(i int) {
			result, found := c.Get(id, d.GetById)

			channel <- Result{result, found}

			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range results {
		results[i] = <-channel
	}

	// Then

	notCachedCount := 0
	validResult := true
	for _, result := range results {
		if !result.cached {
			notCachedCount++
		}		
		if result.data != expectedData {
			validResult = false
		}
	}

	t.Run("should always get valid data from cache and database", func(t* testing.T) {
		if !validResult {
			t.Fatalf("Data from cache is not valid")
		}
	})

	t.Run("should cache data after one request to database and refresh cache after it expires", func(t *testing.T) {
		if notCachedCount < 2 {
			t.Fatalf("Data in cache should be unavailable at least twice, actual: %v", notCachedCount)
		}
	})

}
