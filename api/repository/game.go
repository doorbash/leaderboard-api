package repository

import (
	"context"
	"database/sql"

	"github.com/doorbash/leaderboard-api/api/domain"
)

type GameRepository struct {
	db *sql.DB
}

func (g *GameRepository) GetByID(ctx context.Context, id string) (*domain.Game, error) {
	row := g.db.QueryRowContext(ctx, "SELECT id, name, data FROM games WHERE id = ?", id)
	game := domain.Game{}
	if err := row.Scan(
		&game.ID,
		&game.Name,
		&game.Data,
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
