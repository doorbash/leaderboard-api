package repository

import (
	"context"
	"database/sql"

	"github.com/doorbash/leaderboard-api/api/domain"
)

type LeaderboardDataRepository struct {
	db *sql.DB
}

func (ld *LeaderboardDataRepository) GetByUID(ctx context.Context, uid string) ([]domain.LeaderboardData, error) {
	rows, err := ld.db.QueryContext(ctx, "CALL GET_LEADERBOARD(?)", uid)
	if err != nil {
		return nil, err
	}
	ret := make([]domain.LeaderboardData, 0)
	for rows.Next() {
		ld := domain.LeaderboardData{}
		rows.Scan(&ld.Name, &ld.Value1, &ld.Value2, &ld.Value3)
		ret = append(ret, ld)
	}
	return ret, nil
}

func NewLeaderboardDataRepository(db *sql.DB) *LeaderboardDataRepository {
	return &LeaderboardDataRepository{
		db: db,
	}
}
