package cache

import (
	"sync"
	"time"
)

// Possibly DatabaseRequest should be part of the cache struct instead of the Get function
// This would cause less repetition, but would mean every database request would have a separate cache instance
// It depends on the requirements but wouldnt require big changes so I leave it this way for now
type DatabaseRequest func(string) interface{}

type Cache interface {
	Set(string, interface{})
	Get(string, DatabaseRequest) (data interface{}, cached bool)
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

	c.items[id] = cachedData{
		data:      data,
		expiresAt: time.Now().Add(c.expiresAfter),
	}

	c.mu.Unlock()
}

func (c *cache) Get(id string, databaseRequest DatabaseRequest) (interface{}, bool) {
	c.mu.Lock()

	item, found := c.items[id]

	if !found {
		data := databaseRequest(id)
		c.mu.Unlock()
		c.Set(id, data)
		return data, false
	}

	c.mu.Unlock()

	if item.isExpired() {
		c.deleteExpiredData(id)
		return nil, false
	}

	return item.data, true
}

// TODO: Periodically delete pending data

func (c *cache) deleteExpiredData(id string) {
	c.mu.Lock()
	delete(c.items, id)
	c.mu.Unlock()
	// for key, data := range c.items {
	// 	if data.isExpired() {
	// 		delete(c.items, key)
	// 	}
	// }
}

func (d cachedData) isExpired() bool {
	return d.expiresAt.Before(time.Now())
}
