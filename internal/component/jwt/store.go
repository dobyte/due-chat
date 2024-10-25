package jwt

import (
	"context"
	"github.com/dobyte/due/v2/utils/xconv"
	"github.com/go-redis/redis/v8"
	"time"
)

type store struct {
	redis redis.UniversalClient
}

func (s *store) Get(ctx context.Context, key interface{}) (interface{}, error) {
	return s.redis.Get(ctx, xconv.String(key)).Result()
}

func (s *store) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	return s.redis.Set(ctx, xconv.String(key), value, duration).Err()
}

func (s *store) Remove(ctx context.Context, keys ...interface{}) (value interface{}, err error) {
	list := make([]string, 0, len(keys))
	for _, key := range keys {
		list = append(list, xconv.String(key))
	}

	return s.redis.Del(ctx, list...).Result()
}
