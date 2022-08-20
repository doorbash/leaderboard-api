package domain

import (
	"context"
)

type Player struct {
	ID   int    `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type PlayerRepository interface {
	GetByUID(ctx context.Context, uid string) (*Player, error)
	Insert(ctx context.Context, player *Player) error
	Update(ctx context.Context, player *Player) error
}
