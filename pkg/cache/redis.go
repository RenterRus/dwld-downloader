package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Conn *redis.Client
}

func NewCache(addr string, port int) *Cache {
	return &Cache{
		Conn: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", addr, port),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (c *Cache) Close() {
	c.Conn.Close()
}
