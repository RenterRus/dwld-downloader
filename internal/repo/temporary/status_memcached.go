package temporary

import (
	"context"

	"encoding/json"
	"fmt"
	"time"

	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/sensors"
	"github.com/RenterRus/dwld-downloader/pkg/cache"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	RETOUCH_CACHE  = 120
	TOUCH_LIFETIME = RETOUCH_CACHE + 17
)

type Cache struct {
	c *cache.Cache
	// 			link	files
	links map[string]map[string]struct{}
}

func NewMemCache(conn *cache.Cache) CacheRepo {
	return &Cache{
		c:     conn,
		links: make(map[string]map[string]struct{}),
	}
}

func (c *Cache) GetStatus() (*CacheResponse, error) {
	resp := make(map[string]map[string]TaskResp)

	for _, link := range c.links {
		for file := range link {
			res, err := c.get(file)
			if err != nil {
				return nil, fmt.Errorf("get cache: %w", err)
			}

			var preResp *TaskRequest

			err = json.Unmarshal([]byte(res), &preResp)
			if err != nil {
				return nil, fmt.Errorf("get cache (unmarshal): %w", err)
			}

			if _, ok := resp[preResp.Link]; !ok {
				resp[preResp.Link] = make(map[string]TaskResp)
			}

			resp[preResp.Link][preResp.FileName] = TaskResp{
				MoveTo:       preResp.MoveTo,
				MaxQuality:   preResp.MaxQuality,
				Procentage:   preResp.Procentage,
				Status:       entity.StatusMapping[preResp.Status],
				DownloadSize: preResp.DownloadSize,
				CurrentSize:  preResp.CurrentSize,
				Message:      preResp.Message,
			}
		}
	}

	return &CacheResponse{
		WorkStatus: resp,
		Sensors:    sensors.Sensors(),
	}, nil
}

func (c *Cache) SetStatus(task *TaskRequest) error {
	b, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("SetStatus(Marshal): %w", err)
	}

	err = c.set(task.FileName, string(b), TOUCH_LIFETIME)
	if err != nil {
		return fmt.Errorf("SetStaus(set cache): %w", err)
	}

	if _, ok := c.links[task.Link]; !ok {
		c.links[task.Link] = make(map[string]struct{})
	}
	c.links[task.Link][task.FileName] = struct{}{}

	return nil
}

func (c *Cache) LinkDone(link string) {
	delete(c.links, link)
}

func (c *Cache) Revisor(ctx context.Context) {
	t := time.NewTicker(time.Second * RETOUCH_CACHE)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c.toucer()
		}
	}
}

func (c *Cache) toucer() {
	for _, link := range c.links {
		for file := range link {
			c.c.Conn.Touch(file, TOUCH_LIFETIME)
		}
	}
}

func (c *Cache) set(key, value string, TTLSec int32) error {
	err := c.c.Conn.Set(&memcache.Item{
		Key:        keygen(key),
		Value:      []byte(value),
		Expiration: TTLSec,
	})
	if err != nil {
		return fmt.Errorf("cache set: %w", err)
	}

	return nil
}

func (c *Cache) get(key string) (string, error) {
	item, err := c.c.Conn.Get(keygen(key))
	if err != nil {
		return "", fmt.Errorf("cache get: %w", err)
	}

	return string(item.Value), nil
}
