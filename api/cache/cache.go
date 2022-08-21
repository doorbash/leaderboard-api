package cache

import (
	"context"

	"github.com/doorbash/leaderboard-api/api/util"
)

type Cache interface {
	LoadScripts(ctx context.Context) error
}

func InitCacheScripts(caches ...Cache) error {
	for _, c := range caches {
		ctx, cancel := util.GetContextWithTimeout(context.Background())
		defer cancel()
		if err := c.LoadScripts(ctx); err != nil {
			return err
		}
	}
	return nil
}
