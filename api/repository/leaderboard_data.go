package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doorbash/leaderboard-api/api/domain"
)

var (
	ErrBetterDataAlreadyExists = errors.New("better data already exists")
	ErrLimit                   = errors.New("limit error")
)

type LeaderboardDataRepository struct {
	db *sql.DB
}

func (ld *LeaderboardDataRepository) GetByUID(
	ctx context.Context,
	uid string,
	offset int,
	count int,
) ([]domain.LeaderboardData, error) {
	rows, err := ld.db.QueryContext(ctx, "CALL GET_LEADERBOARD(?, ?, ?)", uid, offset, count)
	if err != nil {
		return nil, err
	}
	ret := make([]domain.LeaderboardData, 0)
	for rows.Next() {
		ld := domain.LeaderboardData{}
		rows.Scan(
			&ld.Name,
			&ld.Value1,
			&ld.Value2,
			&ld.Value3,
		)
		ret = append(ret, ld)
	}
	return ret, nil
}

func (ld *LeaderboardDataRepository) Insert(ctx context.Context, lid string, pid string, v1 int, v2 int, v3 int) error {
	row := ld.db.QueryRowContext(ctx, "CALL NEW_LEADERBOARD_DATA(?, ?, ?, ?, ?)", lid, pid, v1, v2, v3)
	var result int
	if err := row.Scan(&result); err != nil {
		return err
	}
	if result == 2 {
		return ErrBetterDataAlreadyExists
	}
	if result == 3 {
		return ErrLimit
	}
	return nil
}

func NewLeaderboardDataRepository(db *sql.DB) *LeaderboardDataRepository {
	return &LeaderboardDataRepository{
		db: db,
	}
}
