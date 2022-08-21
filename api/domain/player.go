package domain

import (
	"context"
)

type Player struct {
	ID     int    `json:"-"`
	UID    string `json:"uid"`
	Name   string `json:"-"`
	Banned bool   `json:"-"`
}

type PlayerRepository interface {
	GetByUID(ctx context.Context, uid string) (*Player, error)
	Insert(ctx context.Context, player *Player) error
	Update(ctx context.Context, player *Player) error
}
