package cache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	Conn *memcache.Client
}

func NewCache(addr string, port int) *Cache {
	return &Cache{
		Conn: memcache.New(fmt.Sprintf("%s:%d", addr, port)),
	}
}

func (c *Cache) Close() {
	c.Conn.Close()
}
