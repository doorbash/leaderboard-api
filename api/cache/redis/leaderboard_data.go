package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/go-redis/redis/v8"
)

type LeaderboardDataCache struct {
	rdb        *redis.Client
	dataExpiry time.Duration
}

func (c *LeaderboardDataCache) GetByUID(ctx context.Context, uid string) (*string, error) {
	str, err := c.rdb.GetEx(ctx, uid, c.dataExpiry).Result()
	if err != nil {
		return nil, err
	}
	return &str, nil
}

func (c *LeaderboardDataCache) Set(ctx context.Context, uid string, list []domain.LeaderboardData) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	_, err = c.rdb.SetNX(ctx, uid, string(b), c.dataExpiry).Result()
	return err
}

func (c *LeaderboardDataCache) LoadScripts(ctx context.Context) error {
	return nil
}

func NewLeaderboardDataCache(dataExpiry time.Duration) *LeaderboardDataCache {
	rcCache := &LeaderboardDataCache{
		rdb: redis.NewClient(&redis.Options{
			Addr:            REDIS_ADDR,
			Password:        "",
			DB:              REDIS_DATABASE_LEADERBOARDS,
			MaxRetries:      3,
			MinRetryBackoff: REDIS_MIN_RETRY_BACKOFF,
			MaxRetryBackoff: REDIS_MAX_RETRY_BACKOFF,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				log.Println("redis:", "OnConnect()", "Leaderboards")
				return nil
			},
		}),
		dataExpiry: dataExpiry,
	}
	return rcCache
}
