package cache

import (
	"time"
)

type CacheInterface interface {
	Set(string, interface{})
	Get(string) (data interface{}, found bool)
}

type cachedData struct {
	data      interface{}
	expiresAt time.Time
}

type cache struct {
	expiresAfter time.Duration
	items        map[string]cachedData
}

func New(
	expiresAfter time.Duration,
) CacheInterface {
	return &cache{
		expiresAfter: expiresAfter,
		items:        map[string]cachedData{},
	}
}

func (c *cache) Set(id string, data interface{}) {
	c.items[id] = cachedData{
		data:      data,
		expiresAt: time.Now().Add(c.expiresAfter),
	}
	
}

func (c *cache) Get(id string) (interface{}, bool) {
	item, found := c.items[id]
	if !found {
		return nil, false
	}

	if item.isExpired() {
		c.deleteExpiredData(id)
		return nil, false
	}

	return item.data, true
}

// TODO: Periodically delete pending data

func (c *cache) deleteExpiredData(id string) {
	delete(c.items, id)
	// for key, data := range c.items {
	// 	if data.isExpired() {
	// 		delete(c.items, key)
	// 	}
	// }
}

func (d cachedData) isExpired() bool {
	return d.expiresAt.Before(time.Now())
}
