package temporary

import (
	"dwld-downloader/pkg/cache"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	c *cache.Cache
}

func NewMemCache(conn *cache.Cache) *Cache {
	return &Cache{
		c: conn,
	}
}

func (c *Cache) Set(key, value string, TTLSec int32) error {
	err := c.c.Conn.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: TTLSec,
	})
	if err != nil {
		return fmt.Errorf("cache set: %w", err)
	}

	return nil
}

func (c *Cache) Get(key string) (string, error) {
	item, err := c.c.Conn.Get(key)
	if err != nil {
		return "", fmt.Errorf("cache get: %w", err)
	}

	return string(item.Value), nil
}
