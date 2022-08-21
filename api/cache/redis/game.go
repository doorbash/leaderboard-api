package redis

import (
	"context"
	"log"
	"time"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/go-redis/redis/v8"
)

type GameCache struct {
	rdb        *redis.Client
	dataExpiry time.Duration
}

func (c *GameCache) GetByID(ctx context.Context, id string) (*string, error) {
	str, err := c.rdb.GetEx(ctx, id, c.dataExpiry).Result()
	if err != nil {
		return nil, err
	}
	return &str, nil
}

func (c *GameCache) Set(ctx context.Context, game domain.Game) error {
	b, err := game.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = c.rdb.SetNX(ctx, game.ID, string(b), c.dataExpiry).Result()
	return err
}

func (c *GameCache) LoadScripts(ctx context.Context) error {
	return nil
}

func NewGameCache(dataExpiry time.Duration) *GameCache {
	rcCache := &GameCache{
		rdb: redis.NewClient(&redis.Options{
			Addr:            REDIS_ADDR,
			Password:        "",
			DB:              REDIS_DATABASE_GAMES,
			MaxRetries:      3,
			MinRetryBackoff: REDIS_MIN_RETRY_BACKOFF,
			MaxRetryBackoff: REDIS_MAX_RETRY_BACKOFF,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				log.Println("redis:", "OnConnect()", "Games")
				return nil
			},
		}),
		dataExpiry: dataExpiry,
	}
	return rcCache
}
