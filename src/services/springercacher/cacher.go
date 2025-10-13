package springercacher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bdoerfchen/webcmd/src/common/config"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

// A Cacher implementation that uses the memory cache of Victor Springer internally
type springerCacher struct {
	client *cache.Client
}

func New(config *config.CacheConfig) (*springerCacher, error) {
	memAdapter, err := memory.NewAdapter(
		memory.AdapterWithCapacity(int(config.MaxResponsesCached)),
		memory.AdapterWithAlgorithm(memory.LRU), // Least recentlly used response is removed when capacity is reached
	)
	if err != nil {
		return nil, fmt.Errorf("error while creating memory adapter: %w", err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memAdapter),
		cache.ClientWithTTL(time.Duration(config.TTL)),
	)
	if err != nil {
		return nil, fmt.Errorf("error while creating cache client: %w", err)
	}
	return &springerCacher{
		client: cacheClient,
	}, nil
}

func (c *springerCacher) Cache(handler http.HandlerFunc) http.HandlerFunc {
	return c.client.Middleware(handler).(http.HandlerFunc)
}
