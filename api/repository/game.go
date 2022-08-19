package repository

import (
	"context"
	"database/sql"

	"github.com/doorbash/leaderboard-api/api/domain"
)

type GameRepository struct {
	db *sql.DB
}

func (g *GameRepository) GetAll(ctx context.Context) ([]domain.Game, error) {
	rows, err := g.db.QueryContext(ctx, "SELECT id, name FROM games ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	ret := make([]domain.Game, 0)
	for rows.Next() {
		game := domain.Game{}
		rows.Scan(&game.ID, &game.Name)
		ret = append(ret, game)
	}
	return ret, nil
}

func (g *GameRepository) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	row := g.db.QueryRowContext(ctx, "SELECT id, name, version, version_name FROM games WHERE id = ?", id)
	game := domain.Game{}
	if err := row.Scan(
		&game.ID,
		&game.Name,
		&game.Version,
		&game.VersionName,
	); err != nil {
		return nil, err
	}
	return &game, nil
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}
