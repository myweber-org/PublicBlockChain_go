
package cache

import (
	"sync"
	"time"
)

type UserPreference struct {
	Theme     string
	Language  string
	Timezone  string
	UpdatedAt time.Time
}

type CacheEntry struct {
	preferences UserPreference
	expiresAt   time.Time
}

type UserPreferencesCache struct {
	mu       sync.RWMutex
	store    map[string]CacheEntry
	ttl      time.Duration
	maxItems int
}

func NewUserPreferencesCache(ttl time.Duration, maxItems int) *UserPreferencesCache {
	return &UserPreferencesCache{
		store:    make(map[string]CacheEntry),
		ttl:      ttl,
		maxItems: maxItems,
	}
}

func (c *UserPreferencesCache) Set(userID string, prefs UserPreference) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.store) >= c.maxItems {
		c.evictOldest()
	}

	c.store[userID] = CacheEntry{
		preferences: prefs,
		expiresAt:   time.Now().Add(c.ttl),
	}
}

func (c *UserPreferencesCache) Get(userID string) (UserPreference, bool) {
	c.mu.RLock()
	entry, exists := c.store[userID]
	c.mu.RUnlock()

	if !exists {
		return UserPreference{}, false
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.store, userID)
		c.mu.Unlock()
		return UserPreference{}, false
	}

	return entry.preferences, true
}

func (c *UserPreferencesCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.store {
		if oldestKey == "" || entry.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		delete(c.store, oldestKey)
	}
}

func (c *UserPreferencesCache) Remove(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, userID)
}

func (c *UserPreferencesCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.store {
		if now.After(entry.expiresAt) {
			delete(c.store, key)
		}
	}
}