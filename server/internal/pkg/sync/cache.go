package sync

import (
	"time"
)

type cachedItem struct {
	item      any
	expiredAt time.Time
}

type Cache[K comparable] struct {
	cacheMap Map[K, cachedItem]
	// cleanAt  *time.Time
	itemExpiration time.Duration
}

func NewCache[K comparable](itemExpiration time.Duration) Cache[K] {
	return Cache[K]{
		cacheMap:       NewMap[K, cachedItem](),
		itemExpiration: itemExpiration,
	}
}

func (c *Cache[K]) Load(key K) (any, bool) {
	item, ok := c.cacheMap.Load(key)
	if !ok {
		return nil, false
	}

	if item.expiredAt.Before(time.Now()) {
		return nil, false
	}

	return item.item, true
}

func (c *Cache[K]) Store(key K, value any) {
	c.cacheMap.Store(key, cachedItem{
		item:      value,
		expiredAt: time.Now().Add(c.itemExpiration),
	})
}
