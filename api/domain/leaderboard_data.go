package domain

import "context"

type LeaderboardData struct {
	Name   string `json:"name"`
	Value1 int    `json:"value1"`
	Value2 int    `json:"value2"`
	Value3 int    `json:"value3"`
}

type LeaderboardDataRepository interface {
	GetByUID(ctx context.Context, uid string, offset int, count int) ([]LeaderboardData, error)
	GetPlayerPosition(ctx context.Context, lid string, pid string) (int, error)
	Insert(ctx context.Context, lid string, pid string, v1 int, v2 int, v3 int) error
}

type LeaderboardDataCache interface {
	GetByUID(ctx context.Context, uid string) (*string, error)
	Set(ctx context.Context, uid string, list []LeaderboardData) error
}
