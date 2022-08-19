package repository

import (
	"context"
	"database/sql"

	"github.com/doorbash/leaderboard-api/api/domain"
)

type LevelRepository struct {
	db *sql.DB
}

func (l *LevelRepository) GetAllByGameID(ctx context.Context, gid string) ([]domain.Level, error) {
	rows, err := l.db.QueryContext(ctx, "SELECT number, data, min_valid_time, max_valid_score FROM levels WHERE gid = ? ORDER BY number ASC", gid)
	if err != nil {
		return nil, err
	}
	ret := make([]domain.Level, 0)
	for rows.Next() {
		level := domain.Level{
			GameID: gid,
		}
		rows.Scan(&level.Number, &level.Data, &level.MinValidTime, &level.MaxValidScore)
		ret = append(ret, level)
	}
	return ret, nil
}

func (l *LevelRepository) GetByGameIDAndNumber(ctx context.Context, gid string, number int) (*domain.Level, error) {
	row := l.db.QueryRowContext(ctx, "SELECT data, min_valid_time, max_valid_score FROM levels WHERE gid = ? AND number = ?", gid, number)
	level := domain.Level{
		GameID: gid,
		Number: number,
	}
	if err := row.Scan(&level.Data, &level.MinValidTime, &level.MaxValidScore); err != nil {
		return nil, err
	}
	return &level, nil
}

func NewLevelRepository(db *sql.DB) *LevelRepository {
	return &LevelRepository{
		db: db,
	}
}
