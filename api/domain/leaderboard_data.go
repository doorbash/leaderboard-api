package domain

import "context"

type LeaderboardData struct {
	Name   string `json:"name"`
	Value1 int    `json:"value1"`
	Value2 int    `json:"value2"`
	Value3 int    `json:"value3"`
}

type LeaderboardDataRepository interface {
	GetByUID(ctx context.Context, uid string) ([]LeaderboardData, error)
}
