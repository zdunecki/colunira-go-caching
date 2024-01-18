package cache

import (
	"sync"
	"time"
)

const (
	NeverExpires time.Duration = -1
)

// Possibly [DataRequest] should be part of the [cache] struct instead of the [Get] function
// This would cause less repetition, but would mean every database request would have a separate cache instance
// It depends on the requirements but wouldn't require big changes so I leave it this way for now
type DataRequest func(string) interface{}

type Cache interface {
	Set(string, interface{})
	Get(string, DataRequest) (data interface{}, cached bool)
	DeleteExpired()
	Clear()
}

type cachedData struct {
	data      interface{}
	expiresAt time.Time
}

type cache struct {
	expiresAfter time.Duration
	items        map[string]cachedData
	mu           sync.Mutex
}

func New(
	expiresAfter time.Duration,
) Cache {
	return &cache{
		expiresAfter: expiresAfter,
		items:        map[string]cachedData{},
	}
}

func (c *cache) Set(id string, data interface{}) {
	c.mu.Lock()

	var expiresAt time.Time
	if c.expiresAfter == NeverExpires {
		expiresAt = time.Now().AddDate(99, 0, 0)
	} else {
		expiresAt = time.Now().Add(c.expiresAfter)
	}

	c.items[id] = cachedData{
		data:      data,
		expiresAt: expiresAt,
	}

	c.mu.Unlock()
}

func (c *cache) Get(id string, databaseRequest DataRequest) (interface{}, bool) {
	c.mu.Lock()

	item, found := c.items[id]

	if !found {
		data := databaseRequest(id)
		c.mu.Unlock()
		c.Set(id, data)
		return data, false
	}

	c.mu.Unlock()

	if item.expired() {
		c.delete(id)
		return nil, false
	}

	return item.data, true
}

func (c *cache) DeleteExpired() {
	for key, data := range c.items {
		if data.expired() {
			delete(c.items, key)
		}
	}
}

func (c *cache) Clear() {
	c.items = make(map[string]cachedData)
}

func (c *cache) delete(id string) {
	c.mu.Lock()
	delete(c.items, id)
	c.mu.Unlock()

}

func (d cachedData) expired() bool {
	return d.expiresAt.Before(time.Now())
}
