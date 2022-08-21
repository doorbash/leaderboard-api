package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type BannedIPsRedisCache struct {
	rdb        *redis.Client
	dataExpiry time.Duration
}

func (c *BannedIPsRedisCache) IsThisIPBanned(ctx context.Context, ip string) (bool, error) {
	ret, err := c.rdb.Exists(ctx, ip).Result()
	if err != nil {
		return false, err
	}
	return ret == 1, nil
}

func (c *BannedIPsRedisCache) BanThisIP(ctx context.Context, ip string) error {
	_, err := c.rdb.SetNX(ctx, ip, "1", c.dataExpiry).Result()
	return err
}

func (c *BannedIPsRedisCache) LoadScripts(ctx context.Context) error {
	return nil
}

func NewBannedIPsRedisCache(dataExpiry time.Duration) *BannedIPsRedisCache {
	rcCache := &BannedIPsRedisCache{
		rdb: redis.NewClient(&redis.Options{
			Addr:            REDIS_ADDR,
			Password:        "",
			DB:              REDIS_DATABASE_BANNED_IPS,
			MaxRetries:      3,
			MinRetryBackoff: REDIS_MIN_RETRY_BACKOFF,
			MaxRetryBackoff: REDIS_MAX_RETRY_BACKOFF,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				log.Println("redis:", "OnConnect()", "BannedIPs")
				return nil
			},
		}),
		dataExpiry: dataExpiry,
	}
	return rcCache
}
