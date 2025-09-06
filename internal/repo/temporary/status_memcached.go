package temporary

import (
	"context"

	"encoding/json"
	"fmt"
	"time"

	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/sensors"
	"github.com/RenterRus/dwld-downloader/pkg/cache"
)

const (
	RETOUCH_CACHE  = 120
	TOUCH_LIFETIME = RETOUCH_CACHE*2 + 17
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

	for row_link, link := range c.links {
		for file := range link {
			res := c.get(row_link)
			if res == "" {
				fmt.Println("miss cache for", file)
				continue
			}

			var preResp *TaskRequest

			err := json.Unmarshal([]byte(res), &preResp)
			if err != nil {
				fmt.Println("get cache (unmarshal):", err)
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

	c.set(task.Link, string(b))

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
			c.c.Conn.Expire(context.Background(), keygen(file), time.Second*TOUCH_LIFETIME)
		}
	}
}

func (c *Cache) set(key, value string) {
	c.c.Conn.Set(context.Background(), keygen(key), []byte(value), time.Second*TOUCH_LIFETIME)
}

func (c *Cache) get(key string) string {
	item := c.c.Conn.Get(context.Background(), keygen(key))

	res, _ := item.Result()

	return res
}
