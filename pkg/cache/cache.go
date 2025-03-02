package cache

import (
	"sync"
	"time"
)

type WeatherCache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
	expiry  time.Duration
}

type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

func NewWeatherCache(expiry time.Duration) *WeatherCache {
	return &WeatherCache{
		entries: make(map[string]cacheEntry),
		expiry:  expiry,
	}
}

func (c *WeatherCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

func (c *WeatherCache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.expiry),
	}
}
