package repository

import (
	"context"
	"database/sql"

	"github.com/doorbash/leaderboard-api/api/domain"
)

type PlayerRepository struct {
	db *sql.DB
}

func (p *PlayerRepository) GetByUID(ctx context.Context, uid string) (*domain.Player, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, name FROM players WHERE uid = ?", uid)
	player := domain.Player{
		UID: uid,
	}
	if err := row.Scan(&player.ID, &player.Name); err != nil {
		return nil, err
	}
	return &player, nil
}

func (p *PlayerRepository) Insert(ctx context.Context, player *domain.Player) error {
	row := p.db.QueryRowContext(ctx, "CALL NEW_PLAYER(?)", player.Name)
	if err := row.Scan(&player.ID, &player.UID); err != nil {
		return err
	}
	return nil
}

func (p *PlayerRepository) Update(ctx context.Context, player *domain.Player) error {
	_, err := p.db.ExecContext(ctx, "UPDATE players SET name = ? WHERE uid = ?", player.Name, player.UID)
	if err != nil {
		return err
	}
	return nil
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}
